package logging

import "testing"

func TestRotateSize(t *testing.T) {
	SetSizeRotate("./logs", "test", 1024, 7, true)
	for i := 0; i < 10000; i++ {
		Info("i=%d", i)
	}
}
