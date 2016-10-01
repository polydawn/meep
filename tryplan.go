package meep

import (
	"reflect"
)

type TryHandler func(error)

/*
	TryPlan is a declarative error handling plan.

	You can dispatch errors according to several patterns, since both the
	golang stdlib and many libraries have seen fit to use a variety of
	different patterns, and they can't easily be distinguished by type alone
	(e.g. should we typeswitch, or do we need to do actual val/ptr compare):

		TryRoute{ByType: exampleVal error,       Handler: fn}
		TryRoute{ByVal:  ptrOrVal error,         Handler: fn}
		TryRoute{ByFunc: func(error) bool,       Handler: fn}
		TryRoute{CatchAny: true,                 Handler: fn}

	The `By*` fields are used to check whether an error should be handled by
	that route; then the handler is called when a route matches.
	Errors are checked against routes in the order they're listed in your TryPlan.

	Use `ByType` as much as you can.  Meep's typed error helpers should make
	typed errors your default coding style.

	Use `ByVal` where you have to; `io.EOF` is one you must check by value.

	Use `ByFunc` as a last resort -- but if you really have to do something
	complicated, go for it.

	If you want to catch *everything*, set the CatchAny flag.
	If using CatchAny, be sure to add it last -- since it matches any error,
	routes after it will never be called.
*/
type TryPlan []TryRoute

/*
	Checks the TryPlan for handlers that match the error, and invokes the
	first one that does.

	The return value is nil if we found and called a handler, or the
	original error if there was no matching handler.

	If the error parameter was nil, no handler will be called, and the
	result will always be nil.
*/
func (tp TryPlan) Handle(e error) error {
	if e == nil {
		return nil
	}
	for _, tr := range tp {
		if tr.Matches(e) {
			tr.Handler(e)
			return nil
		}
	}
	return e
}

/*
	Like `Handle(e)`, but will panic if the (non-nil) error is not handled.
*/
func (tp TryPlan) MustHandle(e error) {
	e = tp.Handle(e)
	if e != nil {
		panic(e)
	}
}

/*
	A single route, used to compose a TryPlan.

	Typical usage is to set a Handler function, and then treat the other
	fields as if it's a 'union' (that is, only set one of them), like this:

		TryRoute{ByVal: io.EOF, Handler: func(e error) { fmt.Println(e) }}
*/
type TryRoute struct {
	ByType   interface{}
	ByVal    interface{}
	ByFunc   func(error) bool
	CatchAny bool

	Handler TryHandler
}

func (tr TryRoute) Matches(e error) bool {
	if tr.ByType != nil {
		return reflect.TypeOf(e) == reflect.TypeOf(tr.ByType)
	}
	if tr.ByVal != nil {
		return e == tr.ByVal
	}
	if tr.ByFunc != nil {
		return tr.ByFunc(e)
	}
	if tr.CatchAny {
		return true
	}
	return false
}
