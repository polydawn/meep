package meep

import (
	"bytes"
	"io"
)

// Errors with stacks!
type TraitTraceable struct {
	Stack Stack
}

type meepTraceable interface {
	isMeepTraceable() *TraitTraceable
}

func (m *TraitTraceable) isMeepTraceable() *TraitTraceable { return m }

func (m TraitTraceable) IsStackSet() bool {
	return len(m.Stack.Frames) > 0
}

/*
	Return the stack of the error formatted as a human readable string:
	one frame per line.  Each line lists the source file, line number, and
	the name of the function.
*/
func (m TraitTraceable) StackString() string {
	var buf bytes.Buffer
	m.WriteStack(&buf)
	return buf.String()
}

// Same job as StackString; use StackString for convenience, use this for performance.
func (m TraitTraceable) WriteStack(w io.Writer) {
	if len(m.Stack.Frames) == 0 {
		w.Write([]byte("stack info not tracked"))
		return
	}
	for _, fr := range m.Stack.Frames {
		//w.Write(tab)
		w.Write([]byte("·> "))
		w.Write([]byte(fr.String()))
		w.Write(br)
	}
}
