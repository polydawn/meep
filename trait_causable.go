package meep

// Errors with other errors as their cause!
type TraitCausable struct {
	Cause error
}

type meepCausable interface {
	isMeepCausable() *TraitCausable
}

func (m *TraitCausable) isMeepCausable() *TraitCausable { return m }
