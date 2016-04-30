package meep

import "testing"

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
	expect := `Error[meep.Woop]: Wonk="Bonk"; `
	if expect != err.Error() {
		t.Errorf("expected %q, got %q", expect, err.Error())
	}
}
