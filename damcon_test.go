package meep

import (
	"fmt"
	"testing"
)

func TestDamconAbsorbsPanic(t *testing.T) {
	defer DamageControl(func(error) {})
	panic("wat")
}

// Things that don't work: This:
/*
func TestDamconChannelling(t *testing.T) {
	errCh := make(chan error)
	go func() {
		defer func() {
			errCh <- Damcon()
		}()
		panic("wat")
	}()
	<-errCh
}
*/

func TestDamconChannelling(t *testing.T) {
	mockErr := fmt.Errorf("serious probelmz")
	errCh := make(chan error)
	go func() {
		defer DamageControl(func(e error) {
			errCh <- e
		})
		panic(mockErr)
	}()
	err := (<-errCh).(*ErrUnderspecified)
	if err.Cause != mockErr {
		t.Errorf("Drat")
	}
}
