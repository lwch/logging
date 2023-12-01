package logging

import (
	"testing"
)

func TestDebugLog(t *testing.T) {
	Debug("debug")
}

func TestInfoLog(t *testing.T) {
	Info("info")
}

func TestWarningLog(t *testing.T) {
	Warning("warn")
}

func TestErrorLog(t *testing.T) {
	Error("error")
}
