package meep

// Bundles all the behaviors!
type Meep struct {
	TraceableError
	CauseableError
	AutodescribingError
}

// Errors with stacks!
type TraceableError struct {
	Stack Stack
}

// Errors with other errors as their cause!
type CauseableError struct {
	Cause error
}

// Errors that generate their messages automatically from their fields!
type AutodescribingError struct {
	self interface{}
}

// The closest thing you'll get to hierarchical errors: put other errors underneath this one!
type GroupingError struct {
	Specifically error

	// it would be a lot nicer if we could just mark the `specifically` error, taglike.
	// we'd rather that the user's type be able to specify their own bounds!
	// but remember, there's a limit on the utility of this: interfaces declared for faux-hierarchies
	//  don't pay off because you end up casting back to handle things.
	//gather func(error) bool // just a thought.  or `[]error` typeexamples might be better.
}

////

func (m Meep) Error() string {
	return m.AutodescribingError.ErrorMessage() + "\n" + m.TraceableError.StackString()
}
