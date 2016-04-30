package meep

func New(err error) error {
	if m, ok := err.(meepAutodescriber); ok {
		m.isMeepAutodescriber().self = err
	}
	return err
}
