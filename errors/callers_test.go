package errors

import (
	"runtime"
	"testing"
)

func TestCallers(t *testing.T) {
	callers()
	t.Log(runtime.Caller(0))
}

func TestCallers1(t *testing.T) {
	func() {
		callers()
	}()

	pc, file, line, ok := runtime.Caller(0)
	t.Log(file, line, ok)
	t.Log(runtime.FuncForPC(pc).Name())
}
