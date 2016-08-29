package meep

import (
	"reflect"
)

var (
	_ TryHandler = TryHandlerDiscard
)

func TryHandlerDiscard(_ error) {}

func TryHandlerMapto(toTmpl interface{}) TryHandler {
	return tryHandlerMapto{toTmpl}.handle
}

// one of those types that only exists so we can hang a func on it
//  and get stacks that look nice instead of saying "TryHandlerMapto.func1".
type tryHandlerMapto struct{ toTmpl interface{} }

func (h tryHandlerMapto) handle(e error) {
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
