package errors

import (
	"fmt"
	"runtime"
)

type Fields map[string]interface{}

func (fields Fields) Dup() Fields {
	new := make(Fields, len(fields))
	for k, v := range fields {
		new[k] = v
	}

	return new
}

func (fields Fields) String() (res string) {
	for k, v := range fields {
		res += k + ": "
		res += fmt.Sprint(v) + ", "
	}

	if l := len(res); l > 0 {
		res = res[:l-2]
	}
	return res
}

type Builder struct {
	Id         string
	Definition string
	traceStack bool
	pkg        string
	fields     Fields
}

type BuilderOptions func(*Builder)

func WithTraceStack(traced bool) BuilderOptions {
	return func(builder *Builder) {
		builder.traceStack = traced
	}
}

func WithPackage(pkg string) BuilderOptions {
	return func(builder *Builder) {
		builder.pkg = pkg
	}
}

func NewBuilder(Id, Def string, options ...BuilderOptions) *Builder {
	b := &Builder{Id: Id, Definition: Def}
	for _, op := range options {
		op(b)
	}

	if b.traceStack {
		if b.pkg == "" {
			pc, _, _, _ := runtime.Caller(1)
			b.pkg = getPackageName(runtime.FuncForPC(pc).Name())
		}
	}

	b.fields = make(Fields)
	return b
}

func (b Builder) Dump() *Builder {
	b.fields = b.fields.Dup()
	return &b
}

func (b *Builder) WithField(name string, value interface{}) *Builder {
	b.fields[name] = value
	return b
}

func (b *Builder) DumpField(name string, value interface{}) *Builder {
	return b.Dump().WithField(name, value)
}

func (b *Builder) WithFields(fields Fields) *Builder {
	for k, v := range fields {
		b.fields[k] = v
	}

	return b
}

func (b *Builder) DumpFields(fields Fields) *Builder {
	return b.Dump().WithFields(fields)
}

func (b *Builder) New(msg string) error {
	err := &_error{
		Id:     b.Id,
		msg:    msg,
		err:    nil,
		pkg:    b.pkg,
		fields: b.fields.Dup(),
	}

	if b.traceStack {
		err.stacks = callers()
	}

	return err
}

func (b *Builder) Error(args ...interface{}) error {
	return b.New(fmt.Sprint(args...))
}

func (b *Builder) Errorf(format string, args ...interface{}) error {
	return b.New(fmt.Sprintf(format, args...))
}

func (b *Builder) Wrap(err error, args ...interface{}) error {
	err, stacks := b.peelStack(err)
	_err := &_error{
		Id:     b.Id,
		msg:    b.Definition,
		err:    err,
		stacks: stacks,
		pkg:    b.pkg,
		fields: b.fields.Dup(),
	}

	if 0 != len(args) {
		_err.msg = fmt.Sprint(args...)
	}

	if b.traceStack && len(stacks) == 0 {
		_err.stacks = callers()
	}

	return _err
}

func (b *Builder) peelStack(err error) (error, Stack) {
	var stacks Stack
	if err1, ok := err.(*_error); ok {
		stacks = err1.stacks

		// empty stacks.
		err1.stacks = nil
		err = err1
	}

	return err, stacks
}

func (b *Builder) Wrapf(err error, format string, args ...interface{}) error {
	err, stacks := b.peelStack(err)
	_err := &_error{
		Id:     b.Id,
		msg:    b.Definition,
		err:    err,
		stacks: stacks,
		pkg:    b.pkg,
		fields: b.fields.Dup(),
	}

	_err.msg = fmt.Sprintf(format, args...)
	if b.traceStack && len(stacks) == 0 {
		_err.stacks = callers()
	}

	return _err
}

func (b *Builder) ToError() error {
	err := &_error{
		Id:     b.Id,
		msg:    b.Definition,
		err:    nil,
		stacks: nil,
		pkg:    b.pkg,
		fields: b.fields.Dup(),
	}

	if b.traceStack {
		err.stacks = callers()
	}

	return err
}

func (b *Builder) Call(err error, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return b.Wrap(err, args...)
}

func (b *Builder) Callf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return b.Wrapf(err, format, args...)
}
