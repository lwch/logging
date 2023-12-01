package logging

import "testing"

func TestRotateSize(t *testing.T) {
	SetSizeRotate(SizeRotateConfig{
		Level:       LevelInfo,
		Dir:         "./logs",
		Name:        "test",
		Size:        1024,
		Rotate:      7,
		WriteStdout: true,
		WriteFile:   true,
	})
	for i := 0; i < 10000; i++ {
		Info("i=%d", i)
	}
}
