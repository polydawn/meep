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
		handled := TryPlan{{ByType: tr.typeExample, Handler: TryHandlerDiscard}}.Handle(tr.err) == nil
		if handled != tr.shouldMatch {
			t.Errorf("Error %q %s match type example %q", tr.err, negs(tr.shouldMatch), tr.typeExample)
		}
	}
}

func TestTryPredicateVal(t *testing.T) {
	meepVal := &ErrUntypedPanic{}
	tt := []struct {
		err         error
		typeExample error
		shouldMatch bool
	}{
		{meepVal, meepVal, true},
		{&ErrUntypedPanic{}, &ErrUntypedPanic{}, false},
		{&ErrUntypedPanic{}, fmt.Errorf("Cthulhu"), false},
		{fmt.Errorf("lmao"), &ErrUntypedPanic{}, false},
		{fmt.Errorf("lmao"), fmt.Errorf("Cthulhu"), false},
		{fmt.Errorf("lmao"), fmt.Errorf("lmao"), false},
		{io.EOF, io.EOF, true},
	}
	for _, tr := range tt {
		handled := TryPlan{{ByVal: tr.typeExample, Handler: TryHandlerDiscard}}.Handle(tr.err) == nil
		if handled != tr.shouldMatch {
			t.Errorf("Error %q (%p) %s match value %q (%p)", tr.err, tr.err, negs(tr.shouldMatch), tr.typeExample, tr.typeExample)
		}
	}
}

func negs(t bool) string {
	if t {
		return "should"
	}
	return "should not"
}
