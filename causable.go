package meep

type meepCausable interface {
	isMeepCausable() *CausableError
}

func (m *CausableError) isMeepCausable() *CausableError { return m }
