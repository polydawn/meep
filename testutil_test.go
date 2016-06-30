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
