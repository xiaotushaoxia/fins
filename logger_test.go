package fins

import "testing"

func TestLogger(t *testing.T) {
	stdoutLoggerInstance.Printf("aa %d", 1)
}
