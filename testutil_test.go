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
