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
