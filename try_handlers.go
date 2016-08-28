package meep

import (
	"reflect"
)

var (
	_ TryHandler = TryHandlerDiscard
)

func TryHandlerDiscard(_ error) {}

func TryHandlerExplain(toTmpl interface{}) TryHandler {
	return tryHandlerExplain{toTmpl}.handle
}

// one of those types that only exists so we can hang a func on it
//  and get stacks that look nice instead of saying "TryHandlerExplain.func1".
type tryHandlerExplain struct{ toTmpl interface{} }

func (h tryHandlerExplain) handle(e error) {
	typ := reflect.TypeOf(h.toTmpl)
	var err error
	switch typ.Kind() {
	case reflect.Ptr:
		err = reflect.New(typ.Elem()).Interface().(error)
	default:
		err = reflect.Zero(typ).Interface().(error)
	}
	//fmt.Printf("%q  ::  %s: %T\n---\n", err, typ.Kind(), err)
	panic(New(
		err,
		Cause(e),
	))
}
