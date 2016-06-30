package meep

import (
	"os"
	"strings"
)

func stripCwd(in string) string {
	var cwd, _ = os.Getwd()
	return strings.Replace(in, cwd, "", -1)
}

func dropLastNLines(in string, ndrop int) string {
	n := 0
	for i := len(in) - 1; i > 0; i-- {
		if in[i] == '\n' {
			n++
		}
		if n >= ndrop {
			return in[0:i]
		}
	}
	return ""
}

// doesn't preserve whether or not you had a trailing break
// because i don't feel like replacing stdlib `strings.Split` with one that does
func dropLinesContaining(in, drop string) string {
	lines := strings.Split(in, "\n")
	keeps := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, drop) {
			continue
		}
		keeps = append(keeps, line)
	}
	return strings.Join(keeps, "\n")
}

/*
	Ok, this one is really wild.

	Stacks in earlier versions of go would report functions like this...

		meep/fixtures/stack2.go:9: fixtures.func·002

	And newer versions now format them like this:

		meep/fixtures/stack2.go:9: fixtures.BeesBuzz.func1

	(Notice that this only makes an appearance with unnamed funcs.)
	Now, the new version is totally super better, yes, yayy, more contextual
	information by a lightyear.

	But it means our tests either need to alias this part and switch on go versions,
	or... just not look at it.
	This function is to do the latter.
*/
func stripFuncsFromStack(in string) string {
	// NO. THIS IS TERRIBLE
	return ""
}
