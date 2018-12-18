package mlog

import (
	"testing"
)

func TestLogger_P(t *testing.T) {
	l := NewLogger("Test", NoticeLog)
	l.Notice("Test")
	l.Debug("Test")
	l.Info("Test")
	l.Warn("Test")
	l.Error("Test")

	l2 := NewLogger("Example", DebugLog)
	l2.Notice("example")
	l2.Debug("example")
	l2.Info("example")
	l2.Warn("example")
	l2.Error("example")
	ll := l2.GetLogger()
	ll.Println("Test")
}
