package fins

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientClosedError_Unwrap(t *testing.T) {
	var err error
	err = ClientClosedError{}
	assert.True(t, errors.Is(err, net.ErrClosed), "ClientClosedError should be net.ErrClosed")
}

func TestEndCodeError_Unwrap(t *testing.T) {
	err := EndCodeError{0x3001}
	err2 := fmt.Errorf("kk: %w", err)
	var e2 = &EndCodeError{}
	assert.True(t, errors.As(err2, e2), "should as EndCodeError")
	assert.Equal(t, e2.EndCode(), uint16(0x3001), "should be 0x3001")
}
