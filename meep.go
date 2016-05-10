package meep

func New(err error, opts ...Opts) error {
	// for `AutodescribingError`:
	if m, ok := err.(meepAutodescriber); ok {
		m.isMeepAutodescriber().self = err
	}

	// for `TraceableError`:
	if m, ok := err.(meepTraceable); ok {
		m.isMeepTraceable().Stack = *captureStack()
	}

	// for `CausableError`:
	if m, ok := err.(meepCausable); ok {
		for _, o := range opts {
			if o.cause != nil {
				m.isMeepCausable().Cause = o.cause
				break
			}
		}
	}

	// for `GroupingError`:
	//  (nothing really; it's mostly a hint to `AutodescribingError`.)

	return err
}

type Opts struct {
	cause error
}

func Cause(x error) Opts {
	return Opts{cause: New(x)}
}
