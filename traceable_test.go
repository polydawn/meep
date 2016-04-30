package meep_test

import (
	"fmt"

	"."
)

func ExampleTraceableErr() {
	type Woop struct {
		meep.TraceableError
		error
	}
	err := meep.New(&Woop{})
	fmt.Println(err.(*Woop).StackString())

	/// Output:
	// FIXME this is hard to test because of the full local path that pops up :(
}
