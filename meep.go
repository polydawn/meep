package meep

func New(err error) error {
	// for `AutodescribingError`:
	if m, ok := err.(meepAutodescriber); ok {
		m.isMeepAutodescriber().self = err
	}

	// for `TraceableError`:
	if m, ok := err.(meepTraceable); ok {
		m.isMeepTraceable().Stack = *captureStack()
	}

	// for `CauseableError`:
	//  (nothing really; it's mostly a hint to `AutodescribingError`.)

	return err
}
