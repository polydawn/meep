package meep

import "reflect"

type tryMatcher struct {
	predicate func(error) bool
	handler   TryHandler
}

type tryPredicateType struct{ typ reflect.Type }

func (t tryPredicateType) Q(e error) bool {
	return reflect.TypeOf(e) == t.typ
}

type tryPredicateVal struct{ val interface{} }

func (t tryPredicateVal) Q(e error) bool {
	return e == t.val
}

var typeof_ErrUntypedPanic = reflect.TypeOf(&ErrUntypedPanic{})

func trueThunk(error) bool { return true }
