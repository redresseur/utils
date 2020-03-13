package errors

import (
	"fmt"
	"runtime"
)

type Builder struct {
	Id         string
	Definition string
	traceStack bool
	pkg        string
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

	return b
}

func (b *Builder) New(msg string) error {
	err := &_error{
		Id:  b.Id,
		msg: msg,
		err: nil,
		pkg: b.pkg,
	}

	if b.traceStack {
		err.stacks = callers()
	}

	return err
}

func (b *Builder) Wrap(err error, args ...interface{}) error {
	_err := &_error{
		Id:     b.Id,
		msg:    b.Definition,
		err:    err,
		stacks: callers(),
		pkg:    b.pkg,
	}

	if 0 != len(args) {
		_err.msg = fmt.Sprint(args...)
	}

	if b.traceStack {
		_err.stacks = callers()
	}

	return _err
}

func (b *Builder) Wrapf(err error, format string, args ...interface{}) error {
	_err := &_error{
		Id:     b.Id,
		msg:    b.Definition,
		err:    err,
		stacks: callers(),
		pkg:    b.pkg,
	}

	_err.msg = fmt.Sprintf(format, args...)
	if b.traceStack {
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
	}

	if b.traceStack {
		err.stacks = callers()
	}

	return err
}
