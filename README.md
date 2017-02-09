More Expressive Error Patterns
------------------------------

consider the following!

```text
Error[ErrMyApplicationStuck]: Subsystem="blamethis"; TaskNum=42;
    Caused by: Error[ErrNoSpoons]:
        Stack trace:
            ·> /build/path/polydawn/meep/autodescriber_test.go:71: meep.TestAutodescribePlusTraceableCause
            ·> /usr/local/go/src/testing/testing.go:447: testing.tRunner
            ·> /usr/local/go/src/runtime/asm_amd64.s:2232: runtime.goexit
```

That's the output of meep errors...

... where the errors were types:

```golang
type ErrMyApplicationStuck struct {
    meep.TraitCausable
    meep.TraitAutodescribing
    Subsystem string
	TaskNum   int
}
type ErrNoSpoons struct {
    meep.TraitTraceable
    meep.TraitAutodescribing
}
```

... and the error site was:

```golang
err := meep.New(
	&ErrMyApplicationStuck{Subsystem:"blamethis", TaskNum: 42},
	meep.Cause(&ErrNoSpoons{}),
)
```

i dunno if it's entirely obvious why i'm excited, but... this is

- A) **typed errors** with
- B) **near zero keyboard mashing** and
- C) **automatically gorgeous output** including
- D) **stack traces** that
- E) **survive channels** if you want to send the error value to other goroutines and
- F) you can put as many **other fields** in the structs as you want and they're still **tersely prettyprinted**.

PLUS, a whole bunch of easily mixed-in behaviors like capturing a chain of "cause"s to an error -- but only
if you decide your type needs that!

Typed errors are pretty much universally acknowledged to beat the pants off `fmt.Errorf("unmangably handwavey")` stringy errors.
**Now start using them**, because it's *easy*, and you can have stacks and all these other bonuses too!

read more
---------

- ["try"-like handling and dispatch blocks](READMORE.md#try-like-handling-and-dispatch-blocks)
- [Availablity](READMORE.md#availability)
- [Performance](READMORE.md#performance)
