package meep_test

import (
	"strings"
	"testing"

	"."
)

func TestTryHandlerMapto(t *testing.T) {
	type Wonk struct{ meep.AllTraits }
	type Bonk struct{ meep.AllTraits }
	type Tonk struct{ meep.AllTraits }

	var result error
	plan := meep.TryPlan{}.
		Catch(&Wonk{}, meep.TryHandlerMapto(&Bonk{})).
		Catch(&Bonk{}, meep.TryHandlerMapto(&Tonk{})).
		Catch(&Tonk{}, func(e error) {
		actual := e.Error()
		actual = strings.Replace(actual, "\t", "\\t\t", -1)
		actual = strings.Replace(actual, "\n", "\\n\n", -1)
		t.Logf("a fantastic cause tree\n>>>\n%s\n<<<\n", actual)
		result = e
	})

	meep.Try(func() {
		meep.Try(func() {
			meep.Try(func() {
				panic(meep.New(&Wonk{}))
			}, plan)
		}, plan)
	}, plan)

	_ = result.(*Tonk)
	_ = result.(*Tonk).Cause.(*Bonk)
	_ = result.(*Tonk).Cause.(*Bonk).Cause.(*Wonk)
}
