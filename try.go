package meep

func RecoverPanics(fn func()) (e error) {
	defer func() {
		e = coerce(recover())
	}()
	fn()
	return
}

func Try(fn func(), plan TryPlan) {
	plan.MustHandle(RecoverPanics(fn))
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
