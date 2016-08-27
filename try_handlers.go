package meep

var (
	_ TryHandler = TryHandlerDiscard
)

func TryHandlerDiscard(_ error) {}
