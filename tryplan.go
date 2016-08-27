package meep

import (
	"reflect"
)

/*
	TryPlan is an error handling configuration.

	You can dispatch errors according to several patterns, since both the
	golang stdlib and many libraries have seen fit to use a variety of
	different patterns, and they can't easily be distinguished by type alone
	(e.g. should we typeswitch, or do we need to do actual val/ptr compare):

		Catch(kind error, handler)
		CatchVal(ptrOrVal error, handler)
		CatchPredicate(func(error) bool, handler) // use only as a last resort, please
		CatchExotic(handler) // anyone who panics as a non error goes here.  nobody should do that, frankly.
		CatchAll(handler) // called as a last resort

	Errors are checked against handlers in the order they're added.

	If an error passes all the way through a TryPlan without matching any of
	the configured handlers, it is raised as a panic.

	If you want to catch *everything*, use the CatchAll handler.
	If using CatchAll, be sure to add it last -- it's really just syntactic
	sugar for an always-true predicate.
*/
type TryPlan struct {
	matchers []tryMatcher
}

type TryHandler func(error)

func (p TryPlan) Catch(typeExample error, handler TryHandler) TryPlan {
	p.matchers = append(p.matchers, tryMatcher{
		predicate: tryPredicateType{reflect.TypeOf(typeExample)}.Q,
		handler:   handler,
	})
	return p
}

func (p TryPlan) CatchVal(ptrOrVal error, handler TryHandler) TryPlan {
	p.matchers = append(p.matchers, tryMatcher{
		predicate: tryPredicateVal{ptrOrVal}.Q,
		handler:   handler,
	})
	return p
}

func (p TryPlan) CatchPredicate(predicate func(error) bool, handler TryHandler) TryPlan {
	p.matchers = append(p.matchers, tryMatcher{
		predicate: predicate,
		handler:   handler,
	})
	return p
}

func (p TryPlan) CatchExotic(handler TryHandler) TryPlan {
	p.matchers = append(p.matchers, tryMatcher{
		predicate: tryPredicateType{typeof_ErrUntypedPanic}.Q,
		handler:   handler,
	})
	return p
}

func (p TryPlan) CatchAll(handler TryHandler) TryPlan {
	p.matchers = append(p.matchers, tryMatcher{
		predicate: trueThunk,
		handler:   handler,
	})
	return p
}
