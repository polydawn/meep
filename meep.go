package meep

func Meep(err error, opts ...Opts) error {
	return New(err, opts...)
}

func New(err error, opts ...Opts) error {
	opt := flattenOpts(opts)

	// for `TraitAutodescribing`:
	if m, ok := err.(meepAutodescriber); ok {
		m.isMeepAutodescriber().self = err
	}

	// for `TraitTraceable`:
	if m, ok := err.(meepTraceable); ok {
		if !opt.nostack {
			m.isMeepTraceable().Stack = *captureStack()
		}
	}

	// for `TraitCausable`:
	if m, ok := err.(meepCausable); ok {
		if opt.cause != nil {
			m.isMeepCausable().Cause = opt.cause
		}
	}

	return err
}

type Opts struct {
	cause   error
	nostack bool
}

func flattenOpts(opts []Opts) Opts {
	v := Opts{}
	for _, o := range opts {
		if o.cause != nil {
			v.cause = o.cause
		}
		if o.nostack == true {
			v.nostack = true
		}

	}
	return v

}

/*
	Use `Cause` to tell `Meep()` that it should attach another error as a
	cause to the error it's initializating.

	Usage:

		meep.Meep(
			&ErrSomethingCausable{},
			meep.Cause(fmt.Errorf("the root cause")),
		)
*/
func Cause(x error) Opts {
	return Opts{cause: New(x)}
}

/*
	Use `NoStack` to tell `Meep()` that it should skip gathering a stack trace
	for this error, even if it has `TraitTraceable`.

	Usage:

		meep.Meep(
			&ErrUsuallyHasAStacktrace{},
			meep.NoStack(), // skip stacks this time.
		)
*/
func NoStack() Opts {
	return Opts{nostack: true}
}
