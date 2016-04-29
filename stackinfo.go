package meep

import (
	"fmt"
	"path"
	"runtime"
)

type Stack struct {
	Frames []Frame
}

type Frame uintptr

func CaptureStack() *Stack {
	// This looks convoluted (and it is), but there's a reason:
	//  `runtime.Callers` badly wants a uintptr slice, but we want another
	//  type so we can hang e.g. a reasonable stringer method on it.
	// You can't cast alias'd types into each other's slices, so...
	//  there you have it; we're stuck copying.
	var pcs [256]uintptr
	// We offset to skip:
	//  0: runtime.Callers itself
	//  1: this function
	//  2: [start_here]
	n := runtime.Callers(2, pcs[:])
	frames := make([]Frame, n)
	for i := 0; i < n; i++ {
		frames[i] = Frame(pcs[i])
	}
	return &Stack{
		Frames: frames,
	}
}

/*
	`String` returns a human readable form of the frame.

	The string includes the path to the file and the linenumber associated
	with the frame, formatted to match the `file:lineno: ` convention (so
	your IDE, if it supports that convention, may let you click-to-jump);
	following the source location info, the function name is suffixed.
*/
func (pc Frame) String() string {
	if pc == 0 {
		return "unknown:0: unknown"
	}
	pc_actual := uintptr(pc) - 1 // yeah, read `runtime.Callers` *carefully*.
	rtfn := runtime.FuncForPC(pc_actual)
	if rtfn == nil {
		return "unknown:0: unknown"
	}
	file, line := rtfn.FileLine(pc_actual)
	return fmt.Sprintf(
		"%s:%d: %s",
		file,
		line,
		path.Base(rtfn.Name()), // this comes as fq pkg name, so drop "dirs"
	)
}
