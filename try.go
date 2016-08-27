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
*/
type TryPlan struct {
	main     func()
	main2    func() error // we could do this with yet more wrappers but it's preferable to keep stacks shorter where possible
	matchers []tryMatcher
	finally  func()
}

type TryHandler func(error)

/*
	Accepts a function, which will be run with `recover` wrapping, and any
	panics routed according to the `Catch*`'s set up.
*/
func Try(f func()) *TryPlan {
	return &TryPlan{main: f, finally: func() {}}
}

func Try2(f func() error) *TryPlan {
	return &TryPlan{main2: f, finally: func() {}}
}

func (p *TryPlan) Catch(kind error, handler TryHandler) *TryPlan {
	p.matchers = append(p.matchers, tryMatcher{
		predicate: tryPredicateType{reflect.TypeOf(kind)}.Q,
		handler:   handler,
	})
	return p
}

func (p *TryPlan) Finally(f func()) *TryPlan {
	f2 := p.finally
	p.finally = func() {
		f()
		f2()
	}
	return p
}
