// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

const (
	numStackFramesToSkip = 3
	stacktraceDepth      = 32
	unknown              = "unknown"
)

var (
	callers          func() stack
	enableStackTrace = false
)

func init() {
	v, ok := os.LookupEnv("XGO_XERRORS_ENABLE_STACK_TRACE")
	if b, err := strconv.ParseBool(v); ok && err == nil {
		enableStackTrace = b
	}

	EnableStackTrace(enableStackTrace)
}

// EnableStackTrace permits enabling/disabling programmatically the stack trace functionality.
// It is NOT thread-safe.
func EnableStackTrace(enable bool) {
	if !enable {
		callers = func() stack { return nil }
		return
	}

	callers = func() stack {
		var pcs [stacktraceDepth]uintptr
		n := runtime.Callers(numStackFramesToSkip, pcs[:])
		return pcs[0:n]
	}
}

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return unknown
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return unknown
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    source file
//	%d    source line
//	%n    function name
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   function name and path of source file relative to the compile time
//	      GOPATH separated by \n\t (<funcname>\n\t<path>)
//	%+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			fmt.Fprint(s, f.name(), "\n\t", f.file())
		default:
			fmt.Fprint(s, path.Base(f.file()))
		}
	case 'd':
		fmt.Fprint(s, strconv.Itoa(f.line()))
	case 'n':
		fmt.Fprint(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		fmt.Fprint(s, ":")
		f.Format(s, 'd')
	}
}

// MarshalText formats a stacktrace Frame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (f Frame) MarshalText() ([]byte, error) {
	name := f.name()
	if name == unknown {
		return []byte(name), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", name, f.file(), f.line())), nil
}

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each Frame in the stack
//	%v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				fmt.Fprint(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	fmt.Fprint(s, "[")
	for i, f := range st {
		if i > 0 {
			fmt.Fprint(s, " ")
		}
		f.Format(s, verb)
	}
	fmt.Fprint(s, "]")
}

// stack represents a stack of program counters.
type stack []uintptr

// StackTracer is implemented by any value that has a StackTrace method.
// The implementation returns a stack of Frames from innermost (newest) to outermost (oldest).
//
// Stack tracing is an opt-in feature of the package. To do so, either:
// 1) set the environment variable XGO_XERRORS_ENABLE_STACK_TRACE to true, or
// 2) call xerrors.EnableStackTrace(true) programmatically.
type StackTracer interface {
	StackTrace() StackTrace
}

// Format implements the fmt.Formatter interface.
func (s stack) Format(st fmt.State, verb rune) {
	if verb == 'v' && st.Flag('+') {
		for _, pc := range s {
			fmt.Fprintf(st, "\n%+v", Frame(pc))
		}
		return
	}

	if verb == 'v' && st.Flag('#') {
		if s == nil {
			fmt.Fprint(st, "(nil)")
		} else {
			fmt.Fprintf(st, "%v", []uintptr(s))
		}
		return
	}

	format := fmt.Sprintf("%%%c", verb)
	fmt.Fprintf(st, format, []uintptr(s))
}

// StackTrace returns a stack trace.
func (s stack) StackTrace() StackTrace {
	if s == nil {
		return nil
	}

	f := make([]Frame, len(s))
	for i := 0; i < len(f); i++ {
		f[i] = Frame((s)[i])
	}
	return f
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
