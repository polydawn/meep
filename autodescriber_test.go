package meep

import (
	"fmt"
	"os"
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
		CauseableError
		Wonk string
	}
	err := New(&Woop{
		Wonk:           "Bonk",
		CauseableError: CauseableError{fmt.Errorf("lecause")},
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
		CauseableError
		Wonk string
	}
	type Boop struct {
		TraceableError
		AutodescribingError
	}
	err := New(&Woop{
		Wonk: "Bonk",
		CauseableError: CauseableError{
			New(&Boop{}),
		},
	})
	expect := `Error[meep.Woop]: Wonk="Bonk";` + "\n"
	expect += "\t" + `Caused by: Error[meep.Boop]: ` + "\n" // trailing space questionable
	expect += "\t\t" + `Stack trace:` + "\n"
	expect += "\t\t\t" + `·> /autodescriber_test.go:74: meep.TestAutodescribePlusTraceableCause` + "\n"
	expect += "\t\t\t" + `·> /usr/local/go/src/testing/testing.go:447: testing.tRunner` + "\n"
	expect += "\t\t\t" + `·> /usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit` + "\n"
	var cwd, _ = os.Getwd()
	actual := err.Error()
	actual = strings.Replace(actual, cwd, "", -1) // strip the local build path
	if expect != actual {
		t.Errorf("mismatch:\n  expected %q\n       got %q", expect, actual)
	}
	t.Logf("this is what errors with causes that have stacktraces look like :D\n>>>\n%s\n<<<\n", actual)
}
