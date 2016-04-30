package meep

import (
	"bytes"
	"io"
)

type meepTraceable interface {
	isMeepTraceable() *TraceableError
}

func (m *TraceableError) isMeepTraceable() *TraceableError { return m }

/*
	Return the stack of the error formatted as a human readable string:
	one frame per line.  Each line lists the source file, line number, and
	the name of the function.
*/
func (m TraceableError) StackString() string {
	var buf bytes.Buffer
	m.WriteStack(&buf)
	return buf.String()
}

// Same job as StackString; use StackString for convenience, use this for performance.
func (m TraceableError) WriteStack(w io.Writer) {
	if len(m.Stack.Frames) == 0 {
		panic("meep:uninitialized")
	}
	for _, fr := range m.Stack.Frames {
		//w.Write(tab)
		w.Write([]byte("·> "))
		w.Write([]byte(fr.String()))
		w.Write(br)
	}
}
