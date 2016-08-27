package meep

import (
	"fmt"
	"io"
	"testing"
)

func TestTryPredicateType(t *testing.T) {
	tt := []struct {
		err         error
		typeExample error
		shouldMatch bool
	}{
		{&ErrUntypedPanic{}, &ErrUntypedPanic{}, true},
		{&ErrUntypedPanic{}, fmt.Errorf("Cthulhu"), false},
		{fmt.Errorf("lmao"), &ErrUntypedPanic{}, false},
		{fmt.Errorf("lmao"), fmt.Errorf("Cthulhu"), true},
		{io.EOF, io.EOF, true},
		{io.EOF, fmt.Errorf("Cthulhu"), true}, // WOMP!  Though `io.EOF` is meant to be used as a value, it's also an errors.Stringer instance -- watch out!
	}
	for _, tr := range tt {
		res := match(tr.err, TryPlan{}.Catch(tr.typeExample, TryHandlerDiscard)) != nil
		if res != tr.shouldMatch {
			t.Errorf("Error %q %s match type example %q", tr.err, negs(tr.shouldMatch), tr.typeExample)
		}
	}
}

func negs(t bool) string {
	if t {
		return "should"
	}
	return "should not"
}
