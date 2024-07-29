package fins

import (
	"fmt"
	"net"
	"time"
)

// UDPClient errors

type ClientClosedError struct {
}

func (ClientClosedError) Error() string {
	return fmt.Sprintf("error client is closed")
}

func (ClientClosedError) Unwrap() error {
	return net.ErrClosed
}

type ClientClosingError struct {
}

func (ClientClosingError) Error() string {
	return fmt.Sprintf("error client is closing")
}

type EmptyWriteRequestError struct {
}

func (EmptyWriteRequestError) Error() string {
	return fmt.Sprintf("error write request is empty")
}

type ResponseLengthError struct {
	want, got int
}

func (e ResponseLengthError) Error() string {
	return fmt.Sprintf("error response size: want %d, got: %d", e.want, e.got)
}

type ResponseTimeoutError struct {
	duration time.Duration
}

func (e ResponseTimeoutError) Error() string {
	return fmt.Sprintf("Response timeout of %s has been reached", e.duration)
}
func (e ResponseTimeoutError) Timeout() bool   { return true }
func (e ResponseTimeoutError) Temporary() bool { return true }

type IncompatibleMemoryAreaError struct {
	area byte
}

func (e IncompatibleMemoryAreaError) Error() string {
	return fmt.Sprintf("The memory area is incompatible with the data type to be read: 0x%X", e.area)
}

// Driver errors

type BCDBadDigitError struct {
	v   string
	val uint64
}

func (e BCDBadDigitError) Error() string {
	return fmt.Sprintf("Bad digit in BCD decoding: %s = %d", e.v, e.val)
}

type BCDOverflowError struct{}

func (e BCDOverflowError) Error() string {
	return "Overflow occurred in BCD decoding"
}

type EmptyPlcUDPAddress struct{}

func (e EmptyPlcUDPAddress) Error() string {
	return "error empty plc udp address"
}

type EndCodeError struct {
	code uint16
}

func (e EndCodeError) Error() string {
	return fmt.Sprintf("error reported by destination: %s", EndCodeToMsg(e.code))
}

func (e EndCodeError) EndCode() uint16 {
	return e.code
}
