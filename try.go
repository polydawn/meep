package meep

func Try(fn func(), plan TryPlan) {
	defer func() {
		if err := coerce(recover()); err != nil {
			handler := match(err, plan)
			if handler == nil {
				panic(err)
			}
			handler(err)
		}
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
	Runs the `rcvrd` value through the whole `plan`, and
	returns the first TryHandler the plan matches, or nil if no matches.
*/
func match(err error, plan TryPlan) TryHandler {
	// Can you beat this with a case switch?  Absolutely.
	//  We'd have to use reflection to generate one though.
	for _, matcher := range plan.matchers {
		if matcher.predicate(err) {
			return matcher.handler
		}
	}
	return nil
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
