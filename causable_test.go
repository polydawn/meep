package meep_test

import (
	"testing"

	"go.polydawn.net/meep"
)

func TestCausableOpts(t *testing.T) {
	type Woop struct {
		error
		meep.CausableError
		Wonk string
	}
	err := meep.New(
		&Woop{Wonk: "Bonk"},
		meep.Cause(Woop{Wonk: "Tonk"}),
	)
	if err.(*Woop).Wonk != "Bonk" {
		t.Errorf("Bonk somehow became %q", err.(*Woop).Wonk)
	}
	e2 := err.(*Woop).CausableError.Cause
	if e2 == nil {
		t.Errorf("Cause was not initialized")
	}
	if e2.(Woop).Wonk != "Tonk" {
		t.Errorf("Bonk somehow became %q", e2.(Woop).Wonk)
	}
}
