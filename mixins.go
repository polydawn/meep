package meep

// Bundles all the behaviors!
type Meep struct {
	StackableError
	CauseableError
	AutodescribingError
}

// Errors with stacks!
type StackableError struct {
	Frames []uintptr
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
	//gather func(error) bool // just a thought.  or `[]error` typeexamples might be better.
}

////

func (m Meep) Error() string {
	return m.AutodescribingError.ErrorMessage() + "\n" + m.StackableError.StackString()
}

func (m StackableError) StackString() string {
	return "todo stack"
}
