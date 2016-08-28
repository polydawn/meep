package meep_test

import (
	"testing"

	"."
)

func TestTryHandlerExplain(t *testing.T) {
	type Wonk struct{ meep.Meep }
	type Bonk struct{ meep.Meep }
	type Tonk struct{ meep.Meep }

	plan := meep.TryPlan{}.
		Catch(&Wonk{}, meep.TryHandlerExplain(&Bonk{})).
		Catch(&Bonk{}, meep.TryHandlerExplain(&Tonk{})).
		Catch(&Tonk{}, func(e error) { t.Logf("%s", e) })

	meep.Try(func() {
		meep.Try(func() {
			meep.Try(func() {
				panic(meep.New(&Wonk{}))
			}, plan)
		}, plan)
	}, plan)
}
