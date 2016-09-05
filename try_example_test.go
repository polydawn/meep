package meep_test

import (
	"fmt"
	"io"

	"."
)

func ExampleTry() {
	meep.Try(func() {
		panic(meep.New(&meep.AllTraits{}))
	}, meep.TryPlan{
		{ByType: &meep.ErrInvalidParam{},
			Handler: meep.TryHandlerMapto(&meep.ErrProgrammer{})},
		{ByVal: io.EOF,
			Handler: meep.TryHandlerDiscard},
		{CatchAny: true,
			Handler: func(error) {
				fmt.Println("caught wildcard")
			}},
	})

	// Output:
	//	caught wildcard
}
