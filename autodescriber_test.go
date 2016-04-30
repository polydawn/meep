package meep

import (
	"fmt"
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
	expect := `Error[meep.Woop]: Wonk="Bonk";`
	// TODO : // expect += "\n\t" + `Caused by: lecause`
	if expect != err.Error() {
		t.Errorf("expected %q, got %q", expect, err.Error())
	}
}
