package fins

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultResponseTimeoutMillisecond uint = 20 // ms
)

// UDPClient Omron FINS client
// this is concurrent safe
type UDPClient struct {
	localAddr UDPAddress
	plcAddr   UDPAddress
	// config
	responseTimeout  atomic.Int64
	byteOrder        atomic.Value // type: binary.ByteOrder
	readGoroutineNum atomic.Int32

	commLogger

	sid  atomicByte
	resp syncRespSlice

	sf      singleflightOne // avoid Close call twice
	closing atomic.Bool
	wg      sync.WaitGroup

	m      sync.Mutex
	conn   *net.UDPConn
	ctx    context.Context
	cancel context.CancelFunc

	ignoreErrorCode atomic.Value // map[uint16]struct{}{}
}

// NewUDPClient creates a new Omron FINS client
func NewUDPClient(localAddr, plcAddr UDPAddress) (*UDPClient, error) {
	if plcAddr.udpAddress == nil {
		return nil, &net.OpError{Op: "dial", Net: "udp", Err: errors.New("missing address")}
	}
	c := &UDPClient{
		localAddr: localAddr,
		plcAddr:   plcAddr,
	}
	c.SetTimeoutMs(defaultResponseTimeoutMillisecond)
	c.SetReadPacketErrorLogger(&stdoutLogger{})
	c.SetByteOrder(binary.BigEndian)
	c.SetReadGoroutineNum(8)

	c.setConnAndCtx(nil)
	return c, nil
}

// ReadWords Reads words from the PLC data area
func (c *UDPClient) ReadWords(memoryArea byte, address uint16, readCount uint16) ([]uint16, error) {
	readBytes, err := c.ReadBytes(memoryArea, address, readCount)
	if err != nil {
		return nil, err
	}
	return c.bytesToUint16s(readBytes), nil
}

// ReadBytes Reads bytes from the PLC data area
// note: readCount is count of uint16, not count of byte, so len(return) is 2*readCount
func (c *UDPClient) ReadBytes(memoryArea byte, address uint16, readCount uint16) ([]byte, error) {
	return wrapRead(c, func() ([]byte, error) {
		return c.readBytes(memoryArea, address, readCount)
	})
}

// ReadString Reads a string from the PLC data area
// note: readCount is count of uint16, not len of string or count of byte
func (c *UDPClient) ReadString(memoryArea byte, address uint16, readCount uint16) (string, error) {
	data, err := c.ReadBytes(memoryArea, address, readCount)
	if err != nil {
		return "", err
	}
	n := bytes.IndexByte(data, 0)
	if n != -1 {
		data = data[:n]
	}
	return string(data), nil
}

// ReadBits Reads bits from the PLC data area
// note: readCount is count of bool, so len(return) is readCount
func (c *UDPClient) ReadBits(memoryArea byte, address uint16, bitOffset byte, readCount uint16) ([]bool, error) {
	return wrapRead(c, func() ([]bool, error) {
		return c.readBits(memoryArea, address, bitOffset, readCount)
	})
}

// ReadClock Reads the PLC clock
func (c *UDPClient) ReadClock() (t *time.Time, err error) {
	return wrapRead(c, func() (*time.Time, error) {
		r, e := c.sendCommandAndCheckResponse(clockReadCommand())
		if e != nil {
			return nil, e
		}
		return decodeClock(r.data)
	})
}

// WriteWords Writes words to the PLC data area
func (c *UDPClient) WriteWords(memoryArea byte, address uint16, data []uint16) error {
	return c.WriteBytes(memoryArea, address, c.uint16sToBytes(data))
}

// WriteBytes Writes bytes array to the PLC data area
// Example:
//
//	WriteBytes(A, 100, []byte{0x01}) will set A100=256  [01 00]
//	WriteBytes(A, 100, []byte{0x01,0x02}) will set A100=256+2 [01 01]
//	WriteBytes(A, 100, []byte{0x01,0x02,0x01}) will set A100=256+2, A101=256  [01 01 01 00]
//
// Warning:
//
//	if len(b) is not even, I append 0 to the end of b, cause low byte of last memory will be set to 0
//	 A200=1(0x00 0x01), call WriteBytes(A, 100, []byte{0x01}), A200 will be 256(0x01,0x00)
func (c *UDPClient) WriteBytes(memoryArea byte, address uint16, b []byte) error {
	if len(b) == 0 {
		return EmptyWriteRequestError{}
	}
	if len(b)%2 != 0 {
		b = append(b, 0)
	}
	return c.wrapOperate(func() error {
		if err := checkIsWordMemoryArea(memoryArea); err != nil {
			return err
		}
		command := writeCommand(memAddr(memoryArea, address), uint16(len(b)/2), b)
		return c.checkResponse(c.sendCommand(command))
	})
}

// WriteString Writes a string to the PLC data area
// Example:
//
//	WriteString(A, 100, "12") will set A100=[0x31,0x32]
//	WriteString(A, 100, "1") will set A100=[0x31,0x00]
//
// Warning:
//
//	same as WriteBytes, if len([]byte(s)) is not even, I append 0 to the end of b, cause low byte of last memory will be set to 0
func (c *UDPClient) WriteString(memoryArea byte, address uint16, s string) error {
	return c.WriteBytes(memoryArea, address, []byte(s))
}

// WriteBits Writes bits to the PLC data area
// Example:
//
//	WriteBits(A, 100, 0, []bool{true,true}) will set A100=256  [00 03]
//	WriteBits(A, 100, 0, []bool{true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true}) will set A100=65535,A101=1  [FF FF 00 01]
func (c *UDPClient) WriteBits(memoryArea byte, address uint16, bitOffset byte, data []bool) error {
	return c.wrapOperate(func() error {
		if err := checkIsBitMemoryArea(memoryArea); err != nil {
			return err
		}
		l := uint16(len(data))
		bts := make([]byte, 0, l)
		for i := 0; i < int(l); i++ {
			var d byte
			if data[i] {
				d = 0x01
			}
			bts = append(bts, d)
		}
		command := writeCommand(memAddrWithBitOffset(memoryArea, address, bitOffset), l, bts)

		return c.checkResponse(c.sendCommand(command))
	})
}

// SetBit Sets a bit in the PLC data area
func (c *UDPClient) SetBit(memoryArea byte, address uint16, bitOffset byte) error {
	return c.bitTwiddle(memoryArea, address, bitOffset, 0x01)
}

// ResetBit Resets a bit in the PLC data area
// Example:
//
//	ResetBit(A, 100, 0) will set A100.0=0  [00 01] -> [00 00]
//	ResetBit(A, 100, 16) will set A101.0=0  [00 00 00 01] -> [00 00 00 00]
func (c *UDPClient) ResetBit(memoryArea byte, address uint16, bitOffset byte) error {
	return c.bitTwiddle(memoryArea, address, bitOffset, 0x00)
}

// ToggleBit Toggles a bit in the PLC data area
func (c *UDPClient) ToggleBit(memoryArea byte, address uint16, bitOffset byte) error {
	return c.wrapOperate(func() error {
		b, err := c.readBits(memoryArea, address, bitOffset, 1)
		if err != nil {
			return err
		}
		var t byte
		if b[0] {
			t = 0x01
		}
		return c._bitTwiddle(memoryArea, address, bitOffset, t)
	})
}

// SetByteOrder
// Set byte order
// Default value: binary.BigEndian
func (c *UDPClient) SetByteOrder(o binary.ByteOrder) {
	if o != nil {
		c.byteOrder.Store(o)
	}
}

// SetTimeoutMs
// Set response timeout duration (ms).
// Default value: 20ms.
// A timeout of zero can be used to block indefinitely.
func (c *UDPClient) SetTimeoutMs(t uint) {
	c.responseTimeout.Store(int64(time.Duration(t) * time.Millisecond))
}

// SetReadGoroutineNum
// Note: won't stop running goroutine
func (c *UDPClient) SetReadGoroutineNum(count uint8) {
	if count > 0 {
		c.readGoroutineNum.Store(int32(count))
	}
}

// Close Closes an Omron FINS connection
func (c *UDPClient) Close() {
	c.sf.do(c.wrapClose)
}

// SetIgnoreErrorCodes
// Set ignore error codes
func (c *UDPClient) SetIgnoreErrorCodes(codes []uint16) {
	mp := map[uint16]struct{}{}
	for _, code := range codes {
		mp[code] = struct{}{}
	}
	c.ignoreErrorCode.Store(mp)
}

// ============== private ==============
func (c *UDPClient) initConnAndStartReadLoop() error {
	if c.closing.Load() {
		return ClientClosingError{}
	}
	c.m.Lock()
	defer c.m.Unlock()
	if c.conn != nil {
		return nil
	}
	conn, er := net.DialUDP("udp", c.localAddr.udpAddress, c.plcAddr.udpAddress)
	if er != nil {
		return er
	}
	if conn == nil {
		return &net.OpError{Op: "dial", Net: "udp", Err: errors.New("dail return nil conn and nil error")}
	}
	c.setConnAndCtx(conn)

	rn := int(c.readGoroutineNum.Load())
	c.wg.Add(rn)
	for i := 0; i < rn; i++ {
		go func() {
			defer c.wg.Done()
			c.readLoop(c.ctx)
		}()
	}

	return nil
}

func (c *UDPClient) readLoop(ctx context.Context) {
	var buf = make([]byte, udpPacketMaxSize)
	done := ctx.Done()
	for {
		select {
		case <-done:
			return
		default:
			conn := c.getConn()
			if conn == nil { // what happened ?
				c.printFinsPacketError("unknown error: conn is nil in readLoop")
				waitMoment(ctx, time.Millisecond*100)
				continue
			}

			n, _, err := conn.ReadFromUDP(buf)
			if err != nil || n < minResponsePacketSize {
				c.handleReadError(ctx, n, err, buf)
				continue
			}

			respPacket := make([]byte, n)
			copy(respPacket, buf)
			c.printPacket("read", respPacket)
			c.sendToSpecificRespChan(decodeResponse(respPacket))
		}
	}
}

func (c *UDPClient) sendToSpecificRespChan(ans *response) {
	ch := c.resp.getW(ans.header.serviceID)
	if ch == nil {
		c.printFinsPacketError("fins client: no resp chan for sid %d. maybe receive goroutine wait timeout", ans.header.serviceID)
		return
	}

	timeout := time.Duration(c.responseTimeout.Load())
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case ch <- ans:
	case <-c.ctx.Done(): // c.ctx is read only between initConnAndStartReadLoop and Close
		c.printFinsPacketError("fins client: failed to send resp to chan resp. ctx.Done()")
	case <-timer.C:
		c.printFinsPacketError("wait until timeout %s. still no goroutine to receive resp", timeout)
	}
}
func (c *UDPClient) handleReadError(ctx context.Context, n int, err error, buf []byte) {
	if errors.Is(err, net.ErrClosed) {
		return
	}
	msg := fmt.Sprintf("fins client: failed to read fins response packet from: %s: ", c.plcAddr.udpAddress)
	if n < minResponsePacketSize && n > 0 {
		c.printFinsPacketError(msg+"MinResponsePacketSize is %d bytes, got %d bytes: % X", minResponsePacketSize, n, buf[:n])
	} else if n <= 0 {
		c.printFinsPacketError(msg+"ReadFromUDP return %d", n)
	}
	if err != nil {
		c.printFinsPacketError("fins client: failed to ReadFromUDP: " + err.Error())
	}
	waitMoment(ctx, time.Millisecond*100)
	return
}

func (c *UDPClient) createRequest(command []byte) (byte, []byte) {
	sid := c.sid.increment()
	header := defaultCommandHeader(c.localAddr.deviceAddress, c.plcAddr.deviceAddress, sid)
	bts := encodeHeader(header)
	bts = append(bts, command...)
	return sid, bts
}

func (c *UDPClient) sendCommand(command []byte) (*response, error) {
	conn := c.getConn()
	if conn == nil {
		return nil, ClientClosedError{}
	}

	sid, reqPacket := c.createRequest(command)

	respCh := make(chan *response)
	c.resp.set(sid, respCh)
	defer func() {
		c.resp.set(sid, nil)
	}()

	c.printPacket("write", reqPacket)
	_, err := conn.Write(reqPacket)
	if err != nil {
		return nil, err
	}
	var timeoutChan <-chan time.Time
	d := time.Duration(c.responseTimeout.Load())
	if d > 0 {
		timer := time.NewTimer(d)
		defer timer.Stop()
		timeoutChan = timer.C
	}
	select {
	case <-c.ctx.Done(): // c.ctx is read only between initConnAndStartReadLoop and Close
		return nil, ClientClosedError{}
	case respV := <-respCh:
		return respV, nil
	case <-timeoutChan:
		return nil, ResponseTimeoutError{d} // can not actually happen if d == 0
	}
}

func (c *UDPClient) sendCommandAndCheckResponse(command []byte) (*response, error) {
	resp, err := c.sendCommand(command)
	if err != nil {
		return nil, err
	}
	er := c.checkResponse(resp, err)
	if er != nil {
		return nil, er
	}
	return resp, nil
}

func (c *UDPClient) bitTwiddle(memoryArea byte, address uint16, bitOffset byte, value byte) error {
	return c.wrapOperate(func() error {
		if err := checkIsBitMemoryArea(memoryArea); err != nil {
			return err
		}
		return c._bitTwiddle(memoryArea, address, bitOffset, value)
	})
}

func (c *UDPClient) _bitTwiddle(memoryArea byte, address uint16, bitOffset byte, value byte) error {
	mem := memoryAddress{memoryArea, address, bitOffset}
	command := writeCommand(mem, 1, []byte{value})
	return c.checkResponse(c.sendCommand(command))
}

func (c *UDPClient) readBits(memoryArea byte, address uint16, bitOffset byte, readCount uint16) ([]bool, error) {
	if err := checkIsBitMemoryArea(memoryArea); err != nil {
		return nil, err
	}
	command := readCommand(memAddrWithBitOffset(memoryArea, address, bitOffset), readCount)
	r, err := c.sendCommandAndCheckResponse(command)
	if err != nil {
		return nil, err
	}

	result := make([]bool, readCount, readCount)
	for i := 0; i < int(readCount); i++ {
		result[i] = r.data[i]&0x01 > 0
	}
	return result, nil
}
func (c *UDPClient) readBytes(memoryArea byte, address uint16, readCount uint16) ([]byte, error) {
	if err := checkIsWordMemoryArea(memoryArea); err != nil {
		return nil, err
	}
	addr := memAddr(memoryArea, address)
	command := readCommand(addr, readCount)
	r, e := c.sendCommandAndCheckResponse(command)
	if e != nil {
		return nil, e
	}
	if len(r.data) != int(readCount)*2 {
		return nil, ResponseLengthError{want: int(readCount) * 2, got: len(r.data)}
	}
	return r.data, nil
}

func (c *UDPClient) wrapOperate(do func() error) error {
	c.wg.Add(1)
	defer c.wg.Done()
	err := c.initConnAndStartReadLoop()
	if err != nil {
		return err
	}
	return do()
}

func (c *UDPClient) getConn() *net.UDPConn {
	c.m.Lock()
	defer c.m.Unlock()
	return c.conn
}

func (c *UDPClient) setConnAndCtx(conn *net.UDPConn) {
	c.conn = conn
	if conn == nil {
		c.ctx, c.cancel = context.Background(), func() {}
	} else {
		c.ctx, c.cancel = context.WithCancel(context.Background())
	}
}
func (c *UDPClient) closeConn() {
	c.m.Lock()
	defer c.m.Unlock()
	if c.conn == nil {
		return
	}
	c.cancel()
	c.conn.Close()
}

func (c *UDPClient) wrapClose() {
	c.closing.Store(true)
	defer c.closing.Store(false)
	c.closeConn()
	c.wg.Wait()
	// if c.wg.Wait() return, means no goroutine use this client
	// and since c.closing is true,  no new goroutine can use this client
	// so setConnAndCtx can call without protection of c.m
	c.setConnAndCtx(nil)
}

func (c *UDPClient) uint16sToBytes(us []uint16) []byte {
	bts := make([]byte, 2*len(us), 2*len(us))
	order, ok := c.byteOrder.Load().(binary.ByteOrder)
	if !ok {
		order = binary.BigEndian
	}
	for i := 0; i < len(us); i++ {
		order.PutUint16(bts[i*2:i*2+2], us[i])
	}
	return bts
}

func (c *UDPClient) bytesToUint16s(bs []byte) []uint16 {
	order, ok := c.byteOrder.Load().(binary.ByteOrder)
	if !ok {
		order = binary.BigEndian
	}

	data := make([]uint16, len(bs)/2)
	for i := 0; i < len(bs)/2; i++ {
		data[i] = order.Uint16(bs[i*2 : i*2+2])
	}
	return data
}

func (c *UDPClient) checkResponse(r *response, err error) error {
	if err != nil {
		return err
	}
	if r.endCode == EndCodeNormalCompletion {
		return nil
	}
	m, _ := c.ignoreErrorCode.Load().(map[uint16]struct{})
	if _, ok := m[r.endCode]; ok {
		return nil
	}
	return EndCodeError{r.endCode}
}

func wrapRead[T any](c *UDPClient, do func() (T, error)) (result T, err error) {
	c.wg.Add(1)
	defer c.wg.Done()
	if err = c.initConnAndStartReadLoop(); err != nil {
		return result, err
	}
	return do()
}
func decodeClock(data []byte) (*time.Time, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("failed to decode colck: data length should be 6, got: %d", len(data))
	}
	year, err := decodeBCD(data[0:1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode year from %X: %w", data[0:1], err)
	}
	if year < 50 {
		year += 2000
	} else {
		year += 1900
	}
	month, err := decodeBCD(data[1:2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode month from %X: %w", data[1], err)
	}
	day, err := decodeBCD(data[2:3])
	if err != nil {
		return nil, fmt.Errorf("failed to decode day from %X: %w", data[2], err)
	}
	hour, err := decodeBCD(data[3:4])
	if err != nil {
		return nil, fmt.Errorf("failed to decode hour from %X: %w", data[3], err)
	}
	minute, err := decodeBCD(data[4:5])
	if err != nil {
		return nil, fmt.Errorf("failed to decode minute from %X: %w", data[4], err)
	}
	second, err := decodeBCD(data[5:6])
	if err != nil {
		return nil, fmt.Errorf("failed to decode second from % X: %w", data[5:6], err)
	}
	tt := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.Local)
	return &tt, nil
}
