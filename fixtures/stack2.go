package fixtures

// A real mess:
//   BeesBuzz (defers anon) -> beesWuz (defers buzzkill)
//     -> (while returning) busskill:PANIC! -> (while returning) anon -> {fn}.
// ... and then it recovers.
func BeesBuzz(fn func()) {
	defer func() {
		fn()
		recover()
	}(
	// just to see where the defer ends, and if it matters
	)
	beesWuz()
}

func beesWuz() {
	defer buzzkill()
}

func buzzkill() {
	panic("buzznil")
}

// TODO re-panic
// TODO run-time panic (ioob would do)
