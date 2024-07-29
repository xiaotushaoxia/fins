package fins

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

// UDPServer Omron FINS server (PLC emulator)
// it is just for test, only DM area is supported. don't use in production
// fins server is PLC in normal, not our go programs
type UDPServer struct {
	addr      UDPAddress
	conn      *net.UDPConn
	dmarea    []byte
	bitdmarea []byte
	commLogger
	ch chan struct{}
}

const DmAreaSize = 32768

func NewUDPServerSimulator(plcAddr UDPAddress) (*UDPServer, error) {
	if plcAddr.udpAddress == nil { // net.ListenUDP work on random port but I want it fails
		return nil, EmptyPlcUDPAddress{}
	}
	s := new(UDPServer)
	s.addr = plcAddr
	s.dmarea = make([]byte, DmAreaSize)
	s.bitdmarea = make([]byte, DmAreaSize)
	s.ch = make(chan struct{})
	s.SetReadPacketErrorLogger(stdoutLoggerInstance)
	conn, err := net.ListenUDP("udp", plcAddr.udpAddress)
	if err != nil {
		return nil, err
	}
	s.conn = conn

	go func() {
		defer close(s.ch)
		var msg string
		var buf = make([]byte, udpPacketMaxSize) // udp packet max size
		for {
			n, remote, er := conn.ReadFromUDP(buf)
			if er != nil || n < minRequestPacketSize {
				if errors.Is(er, net.ErrClosed) {
					return
				}
				from := ""
				if remote != nil {
					from = fmt.Sprintf(" from: %s", from)
				}
				msg = fmt.Sprintf("fins server %v: failed to read fins request packet%s: ", plcAddr.udpAddress, from)
				if n < minRequestPacketSize && n > 0 {
					s.printFinsPacketError(msg+"minRequestPacketSize is %d, got %d: % X", minResponsePacketSize, n, buf[:n])
				} else if n <= 0 {
					s.printFinsPacketError(msg+"ReadFromUDP return %d", n)
				}
				if err != nil {
					s.printFinsPacketError("fins server %v: failed to ReadFromUDP: %s", plcAddr.udpAddress, err.Error())
				}
				waitMoment(context.Background(), time.Millisecond*100)
				return
			}
			reqPacket := buf[:n]
			s.printPacket("read from "+remote.String(), reqPacket)
			req := decodeRequest(reqPacket)
			resp := s.handler(req)
			respPacket := encodeResponse(resp)
			s.printPacket("write to "+remote.String(), respPacket)
			_, er = conn.WriteToUDP(respPacket, remote)
			if er != nil {
				s.printFinsPacketError("fins server %v: failed to write fins response packet: %s", plcAddr.udpAddress, er)
			}
		}
	}()

	return s, nil
}

// Works with only DM area, 2 byte integers
func (s *UDPServer) handler(r request) response {
	var endCode uint16
	var data []byte
	switch r.commandCode {
	case CommandCodeMemoryAreaRead, CommandCodeMemoryAreaWrite:
		memAddr_ := decodeMemoryAddress(r.data[:4])
		ic := binary.BigEndian.Uint16(r.data[4:6]) // Item count

		switch memAddr_.memoryArea {
		case MemoryAreaDMWord:

			if memAddr_.address+ic*2 > DmAreaSize { // Check address boundary
				endCode = EndCodeAddressRangeExceeded
				break
			}

			if r.commandCode == CommandCodeMemoryAreaRead { //Read command
				data = s.dmarea[memAddr_.address : memAddr_.address+ic*2]
			} else { // Write command
				copy(s.dmarea[memAddr_.address:memAddr_.address+ic*2], r.data[6:6+ic*2])
			}
			endCode = EndCodeNormalCompletion

		case MemoryAreaDMBit:
			if memAddr_.address+ic > DmAreaSize { // Check address boundary
				endCode = EndCodeAddressRangeExceeded
				break
			}
			start := memAddr_.address + uint16(memAddr_.bitOffset)
			if r.commandCode == CommandCodeMemoryAreaRead { //Read command
				data = s.bitdmarea[start : start+ic]
			} else { // Write command
				copy(s.bitdmarea[start:start+ic], r.data[6:6+ic])
			}
			endCode = EndCodeNormalCompletion

		default:
			s.printFinsPacketError("Memory area is not supported: 0x%04x\n", memAddr_.memoryArea)
			endCode = EndCodeNotSupportedByModelVersion
		}

	default:
		s.printFinsPacketError("Command code is not supported: 0x%04x\n", r.commandCode)
		endCode = EndCodeNotSupportedByModelVersion
	}
	return response{defaultResponseHeader(r.header), r.commandCode, endCode, data}
}

// Close Closes the FINS server
func (s *UDPServer) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *UDPServer) Done() <-chan struct{} {
	return s.ch
}
