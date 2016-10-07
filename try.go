package meep

/*
	Invokes your function, captures any panics, and returns them.

	This is simply shorthand for writing your own defers/recovers.

	Note that RecoverPanics returns `error` rather than `interface{}`.
	Any values panicked which do not match the `error` interface will
	be wrapped in the `ErrUntypedPanic` type,
	with the original value in the `Cause` field.
*/
func RecoverPanics(fn func()) (e error) {
	defer func() {
		e = coerce(recover())
	}()
	fn()
	return
}

/*
	Invokes your function, captures any panics, and routes them
	through the TryPlan.

	This may be superior to calling a function that returns an error
	and calling `TryPlan.MustHandle` yourself, because any panics
	that *do* occur will retain a stack including their original panic
	location until the TryPlan evaluates, meaning you can capture it
	properly with any error with the `meep.TraitTraceable` property.
*/
func Try(fn func(), plan TryPlan) {
	// Note that this is subtly different than calling
	//  `plan.MustHandle(RecoverPanics(fn))` ...!
	// This calling MustHandle inside the defer is critical to retaining
	// the stack info in the case of e.g. a nil ptr exception, or anything
	// which doesn't do its own explicit stack capture.
	defer func() {
		plan.MustHandle(coerce(recover()))
	}()
	fn()
}

func coerce(rcvrd interface{}) error {
	switch err := rcvrd.(type) {
	case nil:
		// Panics of nils are possible btw but super absurd.  Never do it.
		return nil
	case error:
		return err
	default:
		// Panics of non-error types are bad and you should feel bad.
		return New(&ErrUntypedPanic{Cause: rcvrd})
	}
}

/*
	A wrapper for non-error types raised from a panic.

	The `Try` system will coerce all non-error types to this automatically.
*/
type ErrUntypedPanic struct {
	TraitAutodescribing
	TraitTraceable
	Cause interface{}
}
