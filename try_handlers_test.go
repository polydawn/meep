package meep_test

import (
	"strings"
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
		Catch(&Tonk{}, func(e error) {
		actual := e.Error()
		actual = strings.Replace(actual, "\t", "\\t\t", -1)
		actual = strings.Replace(actual, "\n", "\\n\n", -1)
		t.Logf("a fantastic cause tree\n>>>\n%s\n<<<\n", actual)
	})

	meep.Try(func() {
		meep.Try(func() {
			meep.Try(func() {
				panic(meep.New(&Wonk{}))
			}, plan)
		}, plan)
	}, plan)
}
