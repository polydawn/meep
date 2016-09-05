package meep

import (
	"os"
	"testing"
)

func BenchmarkTryV1(b *testing.B) {
	var err error
	var val error
	for i := 0; i < b.N; i++ {
		err = &ErrUntypedPanic{}
		Try(
			func() {
				panic(err)
			},
			TryPlan{}.Catch(&ErrUntypedPanic{}, TryHandlerDiscard).
				CatchVal(val, TryHandlerDiscard).
				CatchPredicate(func(e error) bool { return false },
				func(e error) { panic("wow") }).
				CatchAll(TryHandlerDiscard),
		)
	}
}

func BenchmarkTryPlan2Bare(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = &ErrUntypedPanic{}
		TryPlan2{
			{ByType: &ErrUntypedPanic{},
				Handler: TryHandlerDiscard},
			{ByVal: 13,
				Handler: TryHandlerDiscard},
			{ByFunc: func(e error) bool { return false },
				Handler: func(e error) { panic("wow") }},
			{CatchAny: true,
				Handler: TryHandlerDiscard},
		}.Handle(err)
	}
}

func BenchmarkTryPlan2Panicky(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = &ErrUntypedPanic{}
		Try2(
			func() {
				panic(err)
			},
			TryPlan2{
				{ByType: &ErrUntypedPanic{},
					Handler: TryHandlerDiscard},
				{ByVal: 13,
					Handler: TryHandlerDiscard},
				{ByFunc: func(e error) bool { return false },
					Handler: func(e error) { panic("wow") }},
				{CatchAny: true,
					Handler: TryHandlerDiscard},
			},
		)
	}
}

func BenchmarkTypeswitch(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = &ErrUntypedPanic{}
		switch err.(type) {
		case *os.PathError:
			TryHandlerDiscard(err)
		case *os.LinkError:
			panic("wow")
		default:
			TryHandlerDiscard(err)
		}
	}
}

func BenchmarkTypeswitchNofmterr(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = &ErrUntypedPanic{}
		switch err.(type) {
		case *os.PathError:
			TryHandlerDiscard(err)
		case *os.LinkError:
			panic("wow")
		default:
			TryHandlerDiscard(err)
		}
	}
}
