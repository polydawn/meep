package meep_test

import (
	"fmt"

	"."
)

func ExampleTry() {
	meep.Try(func() {
		panic(meep.New(&meep.Meep{}))
	}, meep.TryPlan{}.Catch(&meep.ErrUnderspecified{}, func(error) {
		fmt.Println("caught ErrUnderspecified")
	}).CatchAll(func(error) {
		fmt.Println("caught wildcard")
	}))

	// Output:
	//	caught wildcard
}
