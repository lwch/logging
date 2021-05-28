package logging

import "testing"

func TestLog(t *testing.T) {
	SetDateRotate("./logs", "test", 7, true)
	for i := 0; i < 10000; i++ {
		Info("i=%d", i)
	}
}
