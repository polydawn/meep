package fixtures

// Calls `WheeOne -> wheeTwo -> {fn}`.
func WheeOne(fn func()) {
	wheeTwo(fn)
}

func wheeTwo(fn func()) {
	fn()
}

// Calls `WheeTree (defers wheedee) -> wheeTwo -> (while returning) wheedee -> {fn}`.
func WheeTree(fn func()) {
	defer wheedee(fn)
	wheeTwo(func() {})
}

func wheedee(fn func()) {
	fn()
}
