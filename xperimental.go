package meep

import (
	"io"
	"reflect"
)

func Try2(fn func(), plan TryPlan2) {
	defer func() {
		if err := coerce(recover()); err != nil {
			plan.Handle(err)
			// FIXME fallback should not be silence
			// but `Handle` doesn't report if it *did* anything :I
			// maybe add a `MustHandle`.
		}
	}()
	fn()
}

type TryPlan2 []TryRoute

func (tp TryPlan2) Handle(e error) {
	for _, tr := range tp {
		if tr.Matches(e) {
			tr.Handler(e)
			return
		}
	}
}

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

var plan = TryPlan2{
	{ByType: &ErrUntypedPanic{},
		Handler: TryHandlerDiscard},
	{ByVal: io.EOF,
		Handler: TryHandlerDiscard},
	{ByFunc: func(e error) bool { return true },
		Handler: func(e error) { panic("wow") }},
}
