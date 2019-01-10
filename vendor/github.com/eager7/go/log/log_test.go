package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	l := NewLogger("test", NoticeLog)
	l.Debug("Debug")
	l.Info("Info")
	l.Warn("Warn")
	l.Error("Error")
	//l.Fatal("Fatal")
}



