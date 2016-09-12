package meep

func Meep(err error, opts ...Opts) error {
	return New(err, opts...)
}

func New(err error, opts ...Opts) error {
	// for `TraitAutodescribing`:
	if m, ok := err.(meepAutodescriber); ok {
		m.isMeepAutodescriber().self = err
	}

	// for `TraitTraceable`:
	if m, ok := err.(meepTraceable); ok {
		m.isMeepTraceable().Stack = *captureStack()
	}

	// for `TraitCausable`:
	if m, ok := err.(meepCausable); ok {
		for _, o := range opts {
			if o.cause != nil {
				m.isMeepCausable().Cause = o.cause
				break
			}
		}
	}

	return err
}

type Opts struct {
	cause error
}

func Cause(x error) Opts {
	return Opts{cause: New(x)}
}
