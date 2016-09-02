package meep

import (
	"fmt"
)

/*
	Use `DamageControl` for a safe, terse default mechanism to gather
	panics before they crash your program.

	Damage control can be used like this:

		defer DamageControl(func(e error) {
			errCh <- e
		})

	Using damage control in this way at the top of all your goroutines is
	advised if there's the slightest possibility of panics arising.
	Typically, pushing the error into a channel handled by the spawning
	goroutine is a good response.

	`DamageControl` uses `recover()` internally.  Note that this means you
	must defer `DamageControl` itself (you cannot defer another func which
	calls `DamageControl`; `recover` doesn't work like that).

	The error type given to the `handle` function will always be
	`*ErrUnderspecified`.  If you have different, more specific error
	handling paths and types in mind, you should express those by writing
	your own recovers.
*/
func DamageControl(handle func(error)) {
	if wreck := recover(); wreck != nil {
		err, ok := wreck.(error)
		if !ok {
			err = fmt.Errorf("non-error panic %v", wreck)
		}
		handle(New(&ErrUnderspecified{}, Cause(err)))
	}
}

/*
	A default type for grabbag, underspecified errors; it is the type
	used by `DamageControl` to wrap recovered panics.
*/
type ErrUnderspecified struct {
	TraitAutodescribing
	TraitTraceable
	TraitCausable
}
