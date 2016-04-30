package meep_test

import (
	"fmt"
	"strings"

	"."
)

func ExampleTraceableErr() {
	type Woop struct {
		meep.TraceableError
		error
	}
	err := meep.New(&Woop{})
	str := err.(*Woop).StackString()

	// The *entire* output probably looks something like this:
	//		·> /your/build/path/meep/meep.go:8: meep.New
	//		·> /your/build/path/meep/traceable_test.go:15: meep_test.ExampleTraceableErr
	//		·> /usr/local/go/src/testing/example.go:98: testing.runExample
	//		·> /usr/local/go/src/testing/example.go:36: testing.RunExamplesa
	//		·> /usr/local/go/src/testing/testing.go:486: testing.(*M).Run
	//		·> _/your/build/path/meep/_test/_testmain.go:64: main.main
	//		·> /usr/local/go/src/runtime/proc.go:63: runtime.main
	//		·> /usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit
	// We filter it down rather dramatically so as not to catch any line
	//  numbers from the stdlib we built against, etc.

	str = strings.Split(str, "\n")[1]      // yank the one interesting line
	str = strings.Replace(str, cwd, "", 1) // strip the local build path
	fmt.Println(str)

	// Output:
	//	·> /traceable_test.go:15: meep_test.ExampleTraceableErr
}
