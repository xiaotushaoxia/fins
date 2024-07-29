package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/xiaotushaoxia/fins"
)

var sip *string
var sport, sn, cn, snd, cnd, su, cu, p *int
var client *fins.UDPClient

func init() {
	sip = flag.String("ip", "127.0.0.1", "plc server ip")
	sport = flag.Int("port", 9600, "plc server udp port")

	sn = flag.Int("sn", 0, "plc server network(0-255)")
	cn = flag.Int("cn", 0, "client network(0-255)")

	snd = flag.Int("snd", 0, "plc server node(0-255)")
	cnd = flag.Int("cnd", 0, "client node(0-255)")

	su = flag.Int("su", 0, "plc server unit(0-255)")
	cu = flag.Int("cu", 0, "client unit(0-255)")
	p = flag.Int("p", 1, "show fins udp packet")
	flag.Parse()

	clientAddr := fins.NewUDPAddress("", 0, byte(*cn), byte(*cnd), byte(*cu))
	plcAddr := fins.NewUDPAddress(*sip, *sport, byte(*sn), byte(*snd), byte(*su))
	var err error
	client, err = fins.NewUDPClient(clientAddr, plcAddr)
	if err != nil {
		log.Fatal("failed to NewUDPClient" + err.Error())
	}
	client.SetShowPacket(*p == 1)
}

func main() {
	runCmd(exec)
}

func exec(ss []string) {
	if ss[0] == "set" || ss[0] == "reset" {
		handleSetRest(ss)
		return
	}
	if ss[0] == "close" {
		client.Close()
		return
	}
	if ss[0] == "rc" {
		t, err := client.ReadClock()
		if err != nil {
			fmt.Println("read clock error: " + err.Error())
		} else {
			fmt.Println("read clock success: ", t.Format(time.DateTime))
		}
		return
	}

	rw, mt, dt, addr, count, values, exit := processInputCmdAndShowHelp(ss)
	if exit {
		return
	}
	ma, ok := getMemoryArea(mt, dt)
	if !ok {
		fmt.Println("invalid memory type and data type: ", mt, dt)
		return
	}
	if rw == "r" {
		var result any
		var err error
		switch dt {
		case "b":
			result, err = client.ReadBits(ma, addr, 0, count)
		case "B":
			result, err = client.ReadBytes(ma, addr, count)
		case "s":
			result, err = client.ReadString(ma, addr, count)
		case "w":
			result, err = client.ReadWords(ma, addr, count)
		}
		if err != nil {
			fmt.Println("read error: " + err.Error())
		} else {
			fmt.Println("read success: ", result)
		}
		return
	}

	var err error
	switch dt {
	case "b":
		err = client.WriteBits(ma, addr, 0, values.([]bool))
	case "B":
		err = client.WriteBytes(ma, addr, values.([]byte))
	case "s":
		err = client.WriteString(ma, addr, values.(string))
	case "w":
		err = client.WriteWords(ma, addr, values.([]uint16))
	}
	if err != nil {
		fmt.Println("write error: " + err.Error())
	} else {
		fmt.Println("write success")
	}
}

func handleSetRest(ss []string) {
	if len(ss) != 4 {
		fmt.Println("invalid set/reset input")
		fmt.Println(setResetUsage)
		return
	}
	if ss[1] != "D" && ss[1] != "A" && ss[1] != "H" && ss[1] != "W" {
		fmt.Println(supportMemoryType, "your input: "+ss[1])
		return
	}
	area, ok := getMemoryArea(ss[1], "b")
	if !ok {
		fmt.Println("invalid memory type and data type: ", ss[1], "b")
		return
	}
	addrI, err := strconv.Atoi(ss[2])
	if err != nil || addrI < 0 || addrI > 65535 {
		fmt.Println("invalid address: " + ss[2])
		return
	}
	offsetI, err := strconv.Atoi(ss[3])
	if err != nil || offsetI < 0 || offsetI > 255 {
		fmt.Println("invalid offset: " + ss[3])
		return
	}
	if ss[0] == "set" {
		err = client.SetBit(area, uint16(addrI), byte(offsetI))
	} else {
		err = client.ResetBit(area, uint16(addrI), byte(offsetI))
	}
	if err != nil {
		fmt.Println(ss[0] + " error: " + err.Error())
	} else {
		fmt.Println(ss[0] + " success")
	}
}

func string2bools(s string) ([]bool, error) {
	split := strings.Split(s, ",")
	var bs []bool
	for _, s2 := range split {
		parseBool, err := strconv.ParseBool(s2)
		if err != nil {
			return nil, fmt.Errorf("can not parse %s to []bool: %s", s, s2)
		}
		bs = append(bs, parseBool)
	}
	return bs, nil
}

func string2bytes(s string) ([]byte, error) {
	split := strings.Split(s, ",")
	var bs []byte
	for _, s2 := range split {
		v, err := strconv.Atoi(s2)
		if err != nil || v > 255 || v < 0 {
			return nil, fmt.Errorf("can not parse %s to []byte: %s", s, s2)
		}
		bs = append(bs, byte(v))
	}
	return bs, nil
}

func string2words(s string) ([]uint16, error) {
	split := strings.Split(s, ",")
	var bs []uint16
	for _, s2 := range split {
		v, err := strconv.Atoi(s2)
		if err != nil || v > 65535 || v < 0 {
			return nil, fmt.Errorf("can not parse %s to []word: %s", s, s2)
		}
		bs = append(bs, uint16(v))
	}
	return bs, nil
}

func getMemoryArea(mt, dt string) (byte, bool) {
	switch mt {
	case "D":
		if dt == "b" {
			return fins.MemoryAreaDMBit, true
		}
		return fins.MemoryAreaDMWord, true
	case "A":
		if dt == "b" {
			return fins.MemoryAreaARBit, true
		}
		return fins.MemoryAreaARWord, true
	case "H":
		if dt == "b" {
			return fins.MemoryAreaHRBit, true
		}
		return fins.MemoryAreaHRWord, true
	case "W":
		if dt == "b" {
			return fins.MemoryAreaWRBit, true
		}
		return fins.MemoryAreaWRWord, true
	}
	return 0, false
}

const (
	supportMemoryType = "support memory type: D for DM Area, A for Auxiliary Area, H for Holding Bit Area, W for Work Area"
	supportDataType   = "support data type: b for Bit, B for Byte, s for String, w for Word"
	readUsage         = "read usage:  r <memory type> <data type> <address> <count>  example: r A w 100 1"
	writeUsage        = "write usage: w <memory type> <data type> <address> <values> example: w A w 100 1,2,3"
	setResetUsage     = "set/reset usage: set/reset <memory type> <address> <offset>"
	singleCmdUsage    = "single cmd usage: `close` for close client conn; `rc` for read clock"
)

func help() {
	fmt.Println(supportMemoryType)
	fmt.Println(supportDataType)
	fmt.Println(readUsage)
	fmt.Println(writeUsage)
	fmt.Println(setResetUsage)
	fmt.Println(singleCmdUsage)
}

func processInputCmdAndShowHelp(ss []string) (rw, mt, dt string, addr, count uint16, values any, exit bool) {
	if ss[0] == "h" || ss[0] == "help" {
		help()
		exit = true
		return
	}
	if ss[0] != "r" && ss[0] != "w" {
		fmt.Println("invalid cmd: " + ss[0])
		fmt.Println(readUsage)
		fmt.Println(writeUsage)
		exit = true
		return
	}
	if ss[0] == "r" && len(ss) != 5 {
		fmt.Println("invalid read cmd")
		fmt.Println(readUsage)
		exit = true
		return
	}
	if ss[0] == "w" && len(ss) < 5 {
		fmt.Println("invalid write input")
		fmt.Println(writeUsage)
		exit = true
		return
	}
	if ss[0] == "w" && len(ss) > 5 {
		ss[4] = strings.Join(ss[4:], ",")
		ss = ss[:5]
	}

	if ss[1] != "D" && ss[1] != "A" && ss[1] != "H" && ss[1] != "W" {
		fmt.Println(supportMemoryType, "your input: "+ss[1])
		exit = true
		return
	}

	if ss[2] != "b" && ss[2] != "B" && ss[2] != "s" && ss[2] != "w" {
		fmt.Println(supportDataType, "your input: "+ss[2])
		exit = true
		return
	}

	addrI, err := strconv.Atoi(ss[3])
	if err != nil || addrI < 0 || addrI > 65535 {
		fmt.Println("invalid address: " + ss[3])
		exit = true
		return
	}

	addr = uint16(addrI)
	rw = ss[0]
	mt = ss[1]
	dt = ss[2]
	if rw == "w" {
		valuesStr := ss[4]
		switch dt {
		case "b":
			values, err = string2bools(valuesStr)
		case "B":
			values, err = string2bytes(valuesStr)
		case "s":
			values = valuesStr
		case "w":
			values, err = string2words(valuesStr)
		}
		if err != nil {
			fmt.Println(err)
			exit = true
			return
		}
	} else {
		countI, er := strconv.Atoi(ss[4])
		if er != nil || countI < 0 || countI > 65535 {
			fmt.Println("invalid count: " + ss[4])
			exit = true
			return
		}
		count = uint16(countI)
	}
	return
}

func splitInput(s string) []string {
	ss := strings.Split(s, " ")
	var ns []string
	for _, s2 := range ss {
		if s2 != "" {
			ns = append(ns, s2)
		}
	}
	return ns
}

func runCmd(cmdHandler func(ss []string)) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m>>\033[0m ",
		InterruptPrompt: "^C",
		EOFPrompt:       "cancel",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	var line string
LOOP:
	for {
		line, err = l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		s := strings.TrimSpace(line)

		switch s {
		case "cancel", "bye", "quit", "exit":
			break LOOP
		default:
		}
		input := splitInput(s)
		if len(input) == 0 {
			fmt.Println("empty input")
		} else {
			cmdHandler(input)
		}
	}
}
