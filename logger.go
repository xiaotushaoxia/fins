package fins

import (
	"fmt"
	"sync/atomic"
)

var stdoutLoggerInstance = &stdoutLogger{}

type Logger interface {
	Printf(string, ...any) // auto new line
}

type stdoutLogger struct {
}

func (stdoutLogger) Printf(f string, args ...any) {
	fmt.Println(fmt.Sprintf(f, args...))
}

type commLogger struct {
	readFinsPacketErrLogger atomic.Value
	showPacket              atomic.Bool
}

// SetReadPacketErrorLogger
// read packet run background, we need a Logger to print error
// Default print error to stdout
func (c *commLogger) SetReadPacketErrorLogger(l Logger) {
	c.readFinsPacketErrLogger.Store(l)
}

func (c *commLogger) SetShowPacket(show bool) {
	c.showPacket.Store(show)
}

func (c *commLogger) printFinsPacketError(f string, arg ...any) {
	val := c.readFinsPacketErrLogger.Load()
	if val == nil {
		return
	}
	val.(Logger).Printf(f, arg...)
}

func (c *commLogger) printPacket(rw string, p []byte) {
	if !c.showPacket.Load() {
		return
	}
	val := c.readFinsPacketErrLogger.Load()
	if val == nil {
		return
	}
	val.(Logger).Printf("%s: % X", rw, p)
}
