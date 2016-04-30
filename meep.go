package meep

func New(err error) error {
	if m, ok := err.(meepAutodescriber); ok {
		m.isMeepAutodescriber().self = err
	}
	if m, ok := err.(meepTraceable); ok {
		m.isMeepTraceable().Stack = *CaptureStack()
	}
	return err
}
