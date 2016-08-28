package meep

import (
	"fmt"
	"strings"
	"testing"
)

func TestReacharound(t *testing.T) {
	type Woop struct {
		AutodescribingError
		Wonk string
	}
	var err error
	err = &Woop{Wonk: "Bonk"}
	err = New(err)
	woop := err.(*Woop)
	if woop.Wonk != "Bonk" {
		t.Errorf("Bonk somehow became %q", woop.Wonk)
	}
	if woop.AutodescribingError.self == nil {
		t.Errorf("No impact")
	}
	if woop.AutodescribingError.self != err {
		t.Errorf("Drat")
	}
}

func TestAutodescribeSimple(t *testing.T) {
	type Woop struct {
		AutodescribingError
		Wonk string
	}
	err := New(&Woop{Wonk: "Bonk"})
	expect := `Error[meep.Woop]: Wonk="Bonk";`
	if expect != err.Error() {
		t.Errorf("expected %q, got %q", expect, err.Error())
	}
}

func TestAutodescribePlusCause(t *testing.T) {
	type Woop struct {
		AutodescribingError
		CausableError
		Wonk string
	}
	err := New(&Woop{
		Wonk:          "Bonk",
		CausableError: CausableError{fmt.Errorf("lecause")},
	})
	expect := `Error[meep.Woop]: Wonk="Bonk";` + "\n"
	expect += "\t" + `Caused by: lecause` + "\n"
	actual := err.Error()
	if expect != actual {
		t.Errorf("mismatch:\n  expected %q\n       got %q", expect, actual)
	}
	t.Logf("this is what a very basic error with a nested cause looks like:\n>>>\n%s\n<<<\n", actual)
}

func TestAutodescribePlusTraceableCause(t *testing.T) {
	type Woop struct {
		AutodescribingError
		CausableError
		Wonk string
	}
	type Boop struct {
		TraceableError
		AutodescribingError
	}
	err := New(&Woop{
		Wonk: "Bonk",
		CausableError: CausableError{
			New(&Boop{}),
		},
	})
	expect := `Error[meep.Woop]: Wonk="Bonk";` + "\n"
	expect += "\t" + `Caused by: Error[meep.Boop]:` + "\n"
	expect += "\t\t" + `Stack trace:` + "\n"
	expect += "\t\t\t" + `·> /autodescriber_test.go:73: meep.TestAutodescribePlusTraceableCause` + "\n"

	// Cleanup is fun...
	actual := err.Error()
	// First, remove the local build path for this project.
	actual = stripCwd(actual)
	// Lines we expect following this -- as of go1.4 -- are:
	//   """
	//   expect += "\t\t\t" + `·> /usr/local/go/src/testing/testing.go:447: testing.tRunner` + "\n"
	//   expect += "\t\t\t" + `·> /usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit` + "\n"
	//   """
	// And these are not universal or portable in *several* ways:
	//   - the line numbers aren't constant across go versions
	//   - the files aren't constant across platforms
	//   - indeed even the *number* of lines is not constant across platforms and versions
	//   - the prefix path may change if your GOROOT is unusual (as it is, on some CI platforms, even)
	// So, we must simply truncate them.
	actual = dropLastNLines(actual, 3) + "\n"

	if expect != actual {
		t.Errorf("mismatch:\n  expected %q\n       got %q", expect, actual)
	}
	t.Logf("this is what errors with causes that have stacktraces look like :D\n>>>\n%s\n<<<\n", err.Error())
}

func TestAutodescribePlusTraceableCauseDoubleTrouble(t *testing.T) {
	type Woop struct {
		AutodescribingError
		CausableError
		Wonk string
	}
	type Boop struct {
		AutodescribingError
		CausableError
		TraceableError
	}
	err := New(&Woop{
		Wonk: "Bonk",
		CausableError: CausableError{
			New(&Boop{
				CausableError: CausableError{
					New(&Boop{}),
				},
			}),
		},
	})
	expect := `Error[meep.Woop]: Wonk="Bonk";` + "\n"
	expect += "\t" + `Caused by: Error[meep.Boop]:` + "\n"
	expect += "\t\t" + `Caused by: Error[meep.Boop]:` + "\n"
	expect += "\t\t\t" + `Stack trace:` + "\n"
	expect += "\t\t\t\t" + `·> /autodescriber_test.go:120: meep.TestAutodescribePlusTraceableCauseDoubleTrouble` + "\n"
	// variable: // expect += "\t\t\t\t" + `·> /usr/local/go/src/testing/testing.go:447: testing.tRunner` + "\n"
	// variable: // expect += "\t\t\t\t" + `·> /usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit` + "\n"
	expect += "\t\t" + `Stack trace:` + "\n"
	expect += "\t\t\t" + `·> /autodescriber_test.go:122: meep.TestAutodescribePlusTraceableCauseDoubleTrouble` + "\n"
	// variable: // expect += "\t\t\t" + `·> /usr/local/go/src/testing/testing.go:447: testing.tRunner` + "\n"
	// variable: // expect += "\t\t\t" + `·> /usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit` + "\n"

	// Cleanup is fun...
	actual := err.Error()
	actual = stripCwd(actual)
	actual = dropLinesContaining(actual, ": testing.")
	actual = dropLinesContaining(actual, ": runtime.")

	if expect != actual {
		t.Errorf("mismatch:\n  expected %q\n       got %q", expect, actual)
	}

	// now again for printing (without the parts dropped for the assertion)
	actual = err.Error()
	actual = strings.Replace(actual, "\t", "\\t\t", -1)
	actual = strings.Replace(actual, "\n", "\\n\n", -1)
	t.Logf("this is what errors with causes that have stacktraces look like :D\n>>>\n%s\n<<<\n", actual)
}

func TestAutodescribeManyFields(t *testing.T) {
	type ErrBananaPancakes struct {
		AutodescribingError
		Alpha string
		Beta  int
		Gamma interface{}
		Delta string
	}
	err := New(&ErrBananaPancakes{
		Alpha: "unO",
		Beta:  1,
		Gamma: struct{}{},
		Delta: "ca\ttorce",
	})
	expect := `Error[meep.ErrBananaPancakes]: Alpha="unO";Beta=1;Gamma=struct {}{};Delta="ca\ttorce";`
	actual := err.Error()
	if expect != actual {
		t.Errorf("mismatch:\n  expected %q\n       got %q", expect, actual)
	}
}

func TestIndirectEmbed(t *testing.T) {
	type ErrBananaPancakes struct {
		Meep
		Alpha string
		Beta  int
	}
	err := New(&ErrBananaPancakes{
		Alpha: "fwee",
		Beta:  14,
	}, Cause(fmt.Errorf("a cause")))

	expect := `Error[meep.ErrBananaPancakes]: Alpha="fwee";Beta=14;` + "\n"
	expect += "\t" + `Stack trace:` + "\n"
	expect += "\t\t" + `·> /autodescriber_test.go:184: meep.TestIndirectEmbed` + "\n"
	expect += "\t" + `Caused by: a cause` + "\n"

	// Cleanup is fun...
	actual := err.Error()
	actual = stripCwd(actual)
	actual = dropLinesContaining(actual, ": testing.")
	actual = dropLinesContaining(actual, ": runtime.")
	// Ok, this *really* gets to me.
	// Apparently a multi-line "New" func call:
	//  - in go1.4 and before, it's always seen from the last line.
	//  - in go1.5 and 1.6, is seen from the first line.
	//  - ... except when the race detector is on, in which case it's last line.
	//  - in go1.7 and beyond, it's always seen from the last line again.
	// Wild, right?  Whee.
	// See https://travis-ci.org/polydawn/meep/builds/155778432 .
	// So we'll just take that "first line" line number and uh.  Touch it a little.
	actual = strings.Replace(actual, ":181:", ":184:", 1)

	if expect != actual {
		t.Errorf("mismatch:\n  expected %q\n       got %q", expect, actual)
	}

	// now again for printing (without the parts dropped for the assertion)
	actual = err.Error()
	actual = strings.Replace(actual, "\t", "\\t\t", -1)
	actual = strings.Replace(actual, "\n", "\\n\n", -1)
	t.Logf("embedding meep.Meep shorthand looks about the same\n>>>\n%s\n<<<\n", actual)
}
