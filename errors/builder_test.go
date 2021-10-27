package errors

import (
	"fmt"
	"reflect"
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

func TestBuilder_WithFields(t *testing.T) {
	type fields struct {
		Id         string
		Definition string
		traceStack bool
		pkg        string
		fields     Fields
	}
	type args struct {
		fields Fields
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		{
			name:   "case1",
			fields: fields{fields: Fields{"f1": "xxx", "f2": "yyy"}},
			args: args{
				fields: Fields{"arg1": "zzz", "arg2": "www"},
			},
			want: &Builder{
				fields: Fields{
					"f1": "xxx", "f2": "yyy",
					"arg1": "zzz", "arg2": "www",
				},
			},
		},
		{
			name:   "case2",
			fields: fields{fields: Fields{"f1": "xxx", "f2": "yyy"}},
			args: args{
				fields: Fields{"f1": "zzz", "arg2": "www"},
			},
			want: &Builder{
				fields: Fields{
					"f1":   "zzz",
					"f2":   "yyy",
					"arg2": "www",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				Id:         tt.fields.Id,
				Definition: tt.fields.Definition,
				traceStack: tt.fields.traceStack,
				pkg:        tt.fields.pkg,
				fields:     tt.fields.fields,
			}
			if got := b.WithFields(tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder.WithFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
