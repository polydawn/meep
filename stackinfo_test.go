package meep_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"."
	"./fixtures"
)

var cwd, _ = os.Getwd()

var use14fnnames bool

func init() {
	gover := runtime.Version()
	fmt.Fprintf(os.Stderr, "go version reports as %q\n", gover)
	// I have truely minimal desire to parse this "well".
	// If it's not recognized, we'll assume it's new.
	if gover[0:4] != "go1." {
		return
	}
	switch gover[4] {
	case '0', '1', '2', '3', '4':
		use14fnnames = true
	}
}

type stackFrameExpectation struct {
	n      int
	file   string
	line   int
	func14 string
	func15 string
}

func TestStacksStraightforward(t *testing.T) {
	here := 40
	var result meep.Stack
	fn := func() {
		result = *(meep.CaptureStack())
	}
	fixtures.WheeOne(fn)
	expects := []stackFrameExpectation{
		{0, cwd + "/stackinfo_test.go", here + 3, "meep_test.func·001", "meep_test.TestStacksStraightforward.func1"}, // right here, where we call `CaptureStack`
		{1, cwd + "/fixtures/stack1.go", 9, "", "fixtures.wheeTwo"},                                                  // should be in the body of the func
		{2, cwd + "/fixtures/stack1.go", 5, "", "fixtures.WheeOne"},                                                  // should be in the body of the func
		{3, cwd + "/stackinfo_test.go", here + 5, "", "meep_test.TestStacksStraightforward"},                         // right here, where we call `fixtures.*`
		// No need to get overly precise about line numbers in the stdlib:
		//{4, "/usr/local/go/src/testing/testing.go", 447, "", "testing.tRunner"},
		//{5, "/usr/local/go/src/runtime/asm_amd64.s", 2232, "", "runtime.goexit"},
	}
	expectLen := len(expects) + 2
	assertStack(t, result, expects, expectLen)
}

func TestStacksPlusDeferral(t *testing.T) {
	here := 60
	var result meep.Stack
	fn := func() {
		result = *(meep.CaptureStack())
	}
	fixtures.WheeTree(fn)
	expects := []stackFrameExpectation{
		// note the total lack of 'wheeTwo'; it's called, but already returned before the defer path is hit, so of course it's absent here.
		{0, cwd + "/stackinfo_test.go", here + 3, "meep_test.func·002", "meep_test.TestStacksPlusDeferral.func1"}, // right here, where we call `CaptureStack`
		{1, cwd + "/fixtures/stack1.go", 19, "", "fixtures.wheedee"},                                              // should be in the body of the func (natch, the declare location -- the defer location never shows up; that's not a new func)
		{2, cwd + "/fixtures/stack1.go", 16, "", "fixtures.WheeTree"},                                             // golang considers 'defer' to run on the last line of the parent func.  even if that's "}\n".
		{3, cwd + "/stackinfo_test.go", here + 5, "", "meep_test.TestStacksPlusDeferral"},                         // right here, where we call `fixtures.*`
		// No need to get overly precise about line numbers in the stdlib:
		//{4, "/usr/local/go/src/testing/testing.go", 447, "", "testing.tRunner"},
		//{5, "/usr/local/go/src/runtime/asm_amd64.s", 2232, "", "runtime.goexit"},
	}
	expectLen := len(expects) + 2
	assertStack(t, result, expects, expectLen)
}

func TestStacksPanickingInDefersOhMy(t *testing.T) {
	here := 81
	var result meep.Stack
	fixtures.BeesBuzz(func() {
		result = *(meep.CaptureStack())
	})
	expects := []stackFrameExpectation{
		// note the total lack of reference to where "recover" is called.  (That happened after the stack capture... not that that really matters;
		//   if you flip the recover before the BeesBuzz defer'd func's call to our thunk, this thing on line 9 just moves to 10, that's it -- there's no other flow change.)
		{0, cwd + "/stackinfo_test.go", here + 3, "meep_test.func·003", "meep_test.TestStacksPanickingInDefersOhMy.func1"}, // right here, where we call `CaptureStack` in our thunk
		{1, cwd + "/fixtures/stack2.go", 9, "fixtures.func·002", "fixtures.BeesBuzz.func1"},                                // the line in the deferred function that called our thunk
		// No need to get overly precise about line numbers in the stdlib:
		//{2, "/usr/local/go/src/runtime/asm_amd64.s", 401, "", "runtime.call16"}, // if this isn't a single line on some platforms... uff.
		//{3, "/usr/local/go/src/runtime/panic.go", 387, "", "runtime.gopanic"},  // it might be reasonable to detect these and elide everything following from `runtime.*`.
		{4, cwd + "/fixtures/stack2.go", 22, "", "fixtures.buzzkill"},                              // the line that panicked!
		{5, cwd + "/fixtures/stack2.go", 19, "", "fixtures.beesWuz"},                               // the trailing `}` of `beesWuz`, because we left it via defer
		{6, cwd + "/fixtures/stack2.go", 14, "", "fixtures.BeesBuzz"},                              // the body line the calls down to `beesWuz`
		{7, cwd + "/stackinfo_test.go", here + 4, "", "meep_test.TestStacksPanickingInDefersOhMy"}, // obtw!  when we split the `fixtures.*()` *invocation* across lines, this becomes the last one!
		// No need to get overly precise about line numbers in the stdlib",
		//{8, "/usr/local/go/src/testing/testing.go", 447, "", "testing.tRunner"},
		//{9, "/usr/local/go/src/runtime/asm_amd64.s", 2232, "", "runtime.goexit"},
	}
	expectLen := len(expects) + 4
	assertStack(t, result, expects, expectLen)
}

func assertStack(t *testing.T, result meep.Stack, expects []stackFrameExpectation, expectLen int) {
	// Some quick cleanup on the expectations:
	// If no exceptions were specified, the old 1.4 funcname is expected to be same as the new
	for _, ex := range expects {
		if ex.func14 == "" {
			ex.func14 = ex.func15
		}
	}

	// Assertions!
	for _, tr := range expects {
		file, line, fnname := result.Frames[tr.n].Where()
		if file != tr.file {
			t.Errorf("Stack[%d] file should be %q, was %q", tr.n, tr.file, file)
		}
		if line != tr.line {
			t.Errorf("Stack[%d] line should be %d, was %d", tr.n, tr.line, line)
		}
		expectedFnname := tr.func15
		if use14fnnames {
			expectedFnname = tr.func14
		}
		if fnname != expectedFnname {
			t.Errorf("Stack[%d] func name should be %q, was %q", tr.n, expectedFnname, fnname)
		}
	}
	for i, fr := range result.Frames {
		if i < expectLen {
			continue
		}
		t.Errorf("Stack[%d] was expected to be empty, was %q", i, fr.String())
	}
}
