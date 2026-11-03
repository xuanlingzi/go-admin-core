package logger

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger(WithLevel(TraceLevel), WithName("test"))
	h1 := NewHelper(l).WithFields(map[string]interface{}{"key1": "val1"})
	h1.Trace("trace_msg1")
	h1.Warn("warn_msg1")

	h2 := NewHelper(l).WithFields(map[string]interface{}{"key2": "val2"})
	h2.Trace("trace_msg2")
	h2.Warn("warn_msg2")

	h3 := NewHelper(l).WithFields(map[string]interface{}{"key3": "val4"})
	h3.Info("test_msg")
	ctx := context.TODO()
	ctx = context.WithValue(ctx, &loggerKey{}, h3)
	v := ctx.Value(&loggerKey{})
	ll := v.(*Helper)
	ll.Info("test_msg")
}

// Fields 不能改动共享 logger：请求 A 的字段不能出现在别的请求或直接用 DefaultLogger 打的日志里
func TestFieldsDoNotLeakIntoSharedLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	var mu sync.Mutex
	l := NewLogger(WithLevel(TraceLevel), WithOutput(writerFunc(func(p []byte) (int, error) {
		mu.Lock()
		defer mu.Unlock()
		return buf.Write(p)
	})), WithFields(map[string]interface{}{"svc": "base"}))

	NewHelper(l).WithFields(map[string]interface{}{"x-request-id": "req-A"}).Info("from A")
	l.Log(InfoLevel, "plain sql")
	NewHelper(l).WithFields(map[string]interface{}{"x-request-id": "req-B"}).Info("from B")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got %d: %q", len(lines), buf.String())
	}
	if !strings.Contains(lines[0], "req-A") || !strings.Contains(lines[0], "base") {
		t.Errorf("A line missing its fields: %s", lines[0])
	}
	if strings.Contains(lines[1], "req-A") || !strings.Contains(lines[1], "base") {
		t.Errorf("plain line leaked request fields or lost seed fields: %s", lines[1])
	}
	if strings.Contains(lines[2], "req-A") || !strings.Contains(lines[2], "req-B") {
		t.Errorf("B line wrong: %s", lines[2])
	}
}

type writerFunc func(p []byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }
