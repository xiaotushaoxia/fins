package fins

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFinsClient(t *testing.T) {
	clientAddr := NewUDPAddress("", 9602, 0, 2, 0)
	plcAddr := NewUDPAddress("", 9607, 0, 10, 0)

	c, e := NewUDPClient(clientAddr, plcAddr)
	assert.Nil(t, e)

	for i := 0; i < 10; i++ {
		fmt.Println(i)
		test1(t, c)
		time.Sleep(time.Second)
	}
}
func TestFinsClientRace(t *testing.T) {
	clientAddr := NewUDPAddress("", 9602, 0, 2, 0)
	plcAddr := NewUDPAddress("", 9600, 0, 10, 0)

	c, e := NewUDPClient(clientAddr, plcAddr)
	assert.Nil(t, e)

	for i := 0; i < 10; i++ {
		raceTest(c, i > 5)
	}
}

func test1(t *testing.T, c *UDPClient) {
	defer c.Close()
	s, e := NewUDPServerSimulator(c.plcAddr)
	assert.Nil(t, e)
	defer func() {
		s.Close()
		<-s.Done()
	}()

	toWrite := []uint16{5, 4, 3, 2, 1}

	// ------------- Test Words
	err := c.WriteWords(MemoryAreaDMWord, 100, toWrite)
	assert.Nil(t, err)

	vals, err := c.ReadWords(MemoryAreaDMWord, 100, 5)
	assert.Nil(t, err)
	assert.Equal(t, toWrite, vals)

	// test setting response timeout
	c.SetTimeoutMs(50)
	_, err = c.ReadWords(MemoryAreaDMWord, 100, 5)
	assert.Nil(t, err)

	// ------------- Test Strings
	err = c.WriteString(MemoryAreaDMWord, 10, "ф1234")
	assert.Nil(t, err)

	v, err := c.ReadString(MemoryAreaDMWord, 12, 1)
	assert.Nil(t, err)
	assert.Equal(t, "12", v)

	v, err = c.ReadString(MemoryAreaDMWord, 10, 3)
	assert.Nil(t, err)
	assert.Equal(t, "ф1234", v)

	v, err = c.ReadString(MemoryAreaDMWord, 10, 5)
	assert.Nil(t, err)
	assert.Equal(t, "ф1234", v)

	// ------------- Test Bytes
	err = c.WriteBytes(MemoryAreaDMWord, 10, []byte{0x00, 0x00, 0xC1, 0xA0})
	assert.Nil(t, err)

	b, err := c.ReadBytes(MemoryAreaDMWord, 10, 2)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x00, 0x00, 0xC1, 0xA0}, b)

	buf := make([]byte, 8, 8)
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(-20))
	err = c.WriteBytes(MemoryAreaDMWord, 10, buf)
	assert.Nil(t, err)

	b, err = c.ReadBytes(MemoryAreaDMWord, 10, 4)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x34, 0xc0}, b)

	// ------------- Test Bits
	err = c.WriteBits(MemoryAreaDMBit, 10, 2, []bool{true, false, true})
	assert.Nil(t, err)

	bs, err := c.ReadBits(MemoryAreaDMBit, 10, 2, 3)
	assert.Nil(t, err)
	assert.Equal(t, []bool{true, false, true}, bs)

	bs, err = c.ReadBits(MemoryAreaDMBit, 10, 1, 5)
	assert.Nil(t, err)
	assert.Equal(t, []bool{false, true, false, true, false}, bs)
}

func raceTest(c *UDPClient, testClose bool) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			words, err := c.ReadWords(MemoryAreaDMWord, 100, 10)
			if err != nil {
				fmt.Println(ii, err.Error())
			} else {
				fmt.Println(ii, "read success", words)
			}
		}(i)
	}
	if testClose {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Close()
		}()
	}
	wg.Wait()
}

func TestClient_bytesToUint16s(t *testing.T) {
	var c UDPClient
	v := []uint16{24, 567}
	assert.Equal(t, v, c.bytesToUint16s(c.uint16sToBytes(v)))
}

func Test_atomicValue(t *testing.T) {
	var v atomic.Value
	m, ok := v.Load().(map[uint16]int)
	assert.Equal(t, false, ok)
	assert.Equal(t, 0, m[1])
}
