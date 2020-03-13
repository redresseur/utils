package errors

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
)

/*
	instruction: the errors module
	author: wangzhipengtest@163.com
	date: 2020/03/12
*/

type Stack []runtime.Frame

func (st Stack) format(pkg string) string {
	sts := `[`
	for _, f := range st {
		if pkg == getPackageName(f.Func.Name()) {
			sts += `{file: ` + f.File + ":" + strconv.Itoa(f.Line) + ", func: " + f.Function + "},"
		}
	}
	sts += `]`
	return sts
}

type _error struct {
	Id     string
	msg    string
	err    error `desc:"内置的error"`
	stacks Stack
	pkg    string
}

func (e *_error) Error() string {
	if e.msg != "" {
		return e.msg
	}

	return ""
}

func (e *_error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			msg := "{Id: " + e.Id
			if e.msg != "" {
				msg += ", msg: " + e.msg
			}
			if e.err != nil {
				if _err, ok := e.err.(*_error); ok {
					msg += ", error: "
					io.WriteString(s, msg)
					_err.Format(s, verb)
					msg = ""
				} else {
					msg += ", error: " + fmt.Sprintf("%+v", e.err)
				}
			}
			if 0 != len(e.stacks) {
				msg += ", stack: " + e.stacks.format(e.pkg)
			}

			msg += "}"
			io.WriteString(s, msg)
			return
		}
		fallthrough
	case 's':
		msg := "{Id: " + e.Id
		if e.msg != "" {
			msg += ", msg: " + e.msg
		}
		if e.err != nil {
			if _err, ok := e.err.(*_error); ok {
				msg += ", error: "
				io.WriteString(s, msg)
				_err.Format(s, verb)
				msg = ""
			} else {
				msg += ", error: " + fmt.Sprintf("%v", e.err)
			}
		}

		msg += "}"
		io.WriteString(s, msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

const UNKNOWN_ERROR_ID = `UNKNOWN_ERROR_ID`

func ParseID(err error) string {
	if _err, ok := err.(*_error); ok {
		return _err.Id
	}

	return UNKNOWN_ERROR_ID
}
