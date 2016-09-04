package meep

type ErrInvalidParam struct {
	TraitAutodescribing
	TraitTraceable
	Param  string
	Reason string
}

type ErrNotYetImplemented struct {
	TraitAutodescribing
	TraitTraceable
}

type ErrProgrammer struct {
	TraitAutodescribing
	TraitTraceable
	TraitCausable
}
