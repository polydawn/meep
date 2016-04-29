package meep_test

import (
	"os"
	"testing"

	"."
	"./fixtures"
)

var cwd, _ = os.Getwd()

func TestStacksStraightforward(t *testing.T) {
	var result meep.Stack
	fn := func() {
		result = *(meep.CaptureStack())
	}
	fixtures.WheeOne(fn)
	expect := []struct {
		n   int
		str string
	}{
		{0, cwd + "/stackinfo_test.go:16: meep_test.func·001"},                  // right here, where we call `CaptureStack`
		{1, cwd + "/fixtures/stack1.go:9: fixtures.wheeTwo"},                    // should be in the body of the func
		{2, cwd + "/fixtures/stack1.go:5: fixtures.WheeOne"},                    // should be in the body of the func
		{3, cwd + "/stackinfo_test.go:18: meep_test.TestStacksStraightforward"}, // right here, where we call `fixtures.*`
		// No need to get overly precise about line numbers in the stdlib:
		//{4, "/usr/local/go/src/testing/testing.go:447: testing.tRunner"},
		//{5, "/usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit"},
	}
	expectMax := len(expect) + 2
	for _, tr := range expect {
		str := result.Frames[tr.n].String()
		if str != tr.str {
			t.Errorf("Stack[%d] should be %q, was %q", tr.n, tr.str, str)
		}
	}
	for i, fr := range result.Frames {
		if i < expectMax {
			continue
		}
		t.Errorf("Stack[%d] was expected to be empty, was %q", i, fr.String())
	}
}

func TestStacksPlusDeferral(t *testing.T) {
	var result meep.Stack
	fn := func() {
		result = *(meep.CaptureStack())
	}
	fixtures.WheeTree(fn)
	expect := []struct {
		n   int
		str string
	}{
		// note the total lack of 'wheeTwo'; it's called, but already returned before the defer path is hit, so of course it's absent here.
		{0, cwd + "/stackinfo_test.go:49: meep_test.func·002"},               // right here, where we call `CaptureStack`
		{1, cwd + "/fixtures/stack1.go:19: fixtures.wheedee"},                // should be in the body of the func (natch, the declare location -- the defer location never shows up; that's not a new func)
		{2, cwd + "/fixtures/stack1.go:16: fixtures.WheeTree"},               // golang considers 'defer' to run on the last line of the parent func.  even if that's "}\n".
		{3, cwd + "/stackinfo_test.go:51: meep_test.TestStacksPlusDeferral"}, // right here, where we call `fixtures.*`
		// No need to get overly precise about line numbers in the stdlib:
		//{4, "/usr/local/go/src/testing/testing.go:447: testing.tRunner"},
		//{5, "/usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit"},
	}
	expectMax := len(expect) + 2
	for _, tr := range expect {
		str := result.Frames[tr.n].String()
		if str != tr.str {
			t.Errorf("Stack[%d] should be %q, was %q", tr.n, tr.str, str)
		}
	}
	for i, fr := range result.Frames {
		if i < expectMax {
			continue
		}
		t.Errorf("Stack[%d] was expected to be empty, was %q", i, fr.String())
	}
}
