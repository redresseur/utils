package errors

import (
	"fmt"
	"testing"
)

func TestBuilder_New(t *testing.T) {
	b := NewBuilder("test_err", "this is a test error")
	t.Log(b.New("b.New").Error())
	t.Log(b.New("b.New"))
	t.Logf("%+v", b.New("b.New"))

	b = NewBuilder("test_err", "this is a test error",
		WithTraceStack(true))

	t.Log(b.New("b.New").Error())
	t.Log(b.New("b.New"))
	t.Logf("%+v", b.New("b.New"))
}

func TestBuilder_Wrap(t *testing.T) {
	b := NewBuilder("test_err", "this is a test error")
	t.Log(b.Wrap(b.New("b.Wrap")).Error())
	t.Log(b.Wrap(b.New("b.Wrap")))
	t.Logf("%+v", b.Wrap(b.New("b.Wrap")))

	b = NewBuilder("test_err", "this is a test error",
		WithTraceStack(true))
	t.Log(b.Wrap(b.New("b.Wrap")).Error())
	t.Log(b.Wrap(b.New("b.Wrap")))
	t.Logf("%+v", b.Wrap(b.New("b.Wrap")))
}

func TestBuilder_Wrapf(t *testing.T) {
	b := NewBuilder("test_err", "this is a test error")
	t.Log(b.Wrapf(b.New("b.Wrapf"), "test_wrapf").Error())
	t.Log(b.Wrapf(b.New("b.Wrapf"), "test_wrapf"))
	t.Logf("%+v", b.Wrapf(b.New("b.Wrapf"), "test_wrapf"))

	b = NewBuilder("test_err", "this is a test error",
		WithTraceStack(true))
	t.Log(b.Wrapf(b.New("b.Wrapf"), "test_wrapf").Error())
	t.Log(b.Wrapf(b.New("b.Wrapf"), "test_wrapf"))
	t.Logf("%+v", b.Wrapf(b.New("b.Wrapf"), "test_wrapf"))
}

func BenchmarkBuilder_New(b *testing.B) {
	builder := NewBuilder("test_err", "this is a test error", WithTraceStack(true))
	err := builder.New("this is a benchmark test")
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%+v", err)
	}
}
