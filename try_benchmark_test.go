package meep

import (
	"os"
	"testing"
)

func BenchmarkTryPlanBare(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = &ErrUntypedPanic{}
		TryPlan{
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

func BenchmarkTryPlanPanicky(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = &ErrUntypedPanic{}
		Try(
			func() {
				panic(err)
			},
			TryPlan{
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

func BenchmarkStackinitialization(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = New(&ErrUntypedPanic{})
	}
	_ = err
}

func BenchmarkBaseline(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				recover()
			}()
			func() {
				err = &ErrUntypedPanic{}
			}()
		}()
	}
}
