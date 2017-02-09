A bigger example
----------------

Suppose you have an application which speaks a json protocol.
The basic errors returned by the json.Marshal/Unmarshal functions are great;
but in the context of your larger application, don't tell you things like
"what kind of message was this that caused an error?" -- you need to attach your
own information about that.  So, custom error types to the rescue!

```golang
// Our message type.  Has polymorphic content.
type Envelope struct {
	MsgType  string
	Msg      interface{}
}
type AppleMsg struct{ Opacity int }
type PearMsg struct{ /*...*/ }
```

We have a couple different message types, and one "envelope" type at the top to contain
our protocol header and some hint for what the body message is.
(This kind of polymorphism for json decoding probably looks familiar already; if not,
we're loosely riffing off [eagain's dynamic json tutorial](http://eagain.net/articles/go-dynamic-json/), which is a great document.)

So when decoding this protocol, there's a couple of things that could go wrong:

1. we could fail to parse the json at all;
2. we could see a `MsgType` we don't understand;
3. and we could fail while mapping the inner `Msg` into our structs.

```golang
type ErrBadProtocolHandshake struct {
    meep.TraitCausable
    meep.TraitAutodescribing
}
type ErrMalformedMessage struct {
    meep.TraitCausable
    meep.TraitAutodescribing
	ExpectedType string
}
```

Here we've declared two types of errors.
We're going to lump those first three failure modes into `ErrBadProtocolHandshake`,
and the fourth we'll call `ErrMalformedMessage`.

Notice that `ErrMalformedMessage` has an additional field in it:
when that error comes up, we'll be in a conditional branch
depending on which kind of message body the envelope declared,
so we'll want to attach that information.

Now, let's parse!

Our example message will be:

```golang
envelopeRaw := []byte(`{"MsgType":"apple", "Msg":{"Opacity":"stringy"}}`)
```

(You can see we're lining up for an error because strings aren't ints, here, on the inner message.)

Our parse logic, complete with all error returns, looks like this:

```golang
func() (*Envelope, error) {
	msgRaw := json.RawMessage{}
	msgEnvelope := &Envelope{Msg: &msgRaw}
	if err := json.Unmarshal(envelopeRaw, msgEnvelope); err != nil {
		return nil, meep.New(
			&ErrBadProtocolHandshake{},
			meep.Cause(err),
		)
	}
	var msg interface{}
	switch msgEnvelope.MsgType {
	case "apple":
		msg = &AppleMsg{}
	case "pear":
		msg = &PearMsg{}
	default:
		return msgEnvelope, meep.New(
			&ErrBadProtocolHandshake{},
			meep.Cause(fmt.Errorf("unknown message type")),
		)
	}
	if err := json.Unmarshal(msgRaw, msg); err != nil {
		return msgEnvelope, meep.New(
			&ErrMalformedMessage{ExpectedType: msgEnvelope.MsgType},
			meep.Cause(err),
		)
	}
	msgEnvelope.Msg = msg
	return msgEnvelope, nil
}
```

Now, when an error is returned, you can handle it with a type switch:

```golang
switch err.(type) {
	case *ErrBadProtocolHandshake: /* ... */
	case *ErrMalformedMessage: /* ... */
}
```

And if we print it?

```text
Error[meep_test.ErrMalformedMessage]: ExpectedType="apple";
	Caused by: json: cannot unmarshal string into Go value of type int
```

The error type is clearly displayed.
Any additional fields (in this case, just "ExpectedType") are printed along with it.
The cause (and recursively, if a meep error is the cause of another meep error) is printed on the next line, indented and clearly separated.

Notice that we didn't opt-in to stack traces here -- none of our error types embedded `TraitTraceable`.
In this example, forgoing stack traces seemed reasonable because our error types should be expressive enough;
but you can add stack traces to any error by embedding that trait, and they'll be automatically pretty-printed.
If you do opt-in to stacks by embedding that trait, they're formatted like this:

```
Error[ErrTypeName]:
	Stack trace:
		·> /build/path/polydawn/meep/trait_autodescriber_test.go:120: meep.TestAutodescribePlusTraceableCauseDoubleTrouble
		·> /usr/local/go/src/testing/testing.go:610: testing.tRunner
		·> /usr/local/go/src/runtime/asm_amd64.s:2086: runtime.goexit
```

Everything you did here, you could have done with error type declarations before `meep` came along.
However, to get the same level of features, you would need to write custom stringer methods for each error type,
come up with some sort of a solution for stack traces, and just generally end up with more than 2 or 3x the amount of code.
With `meep`, it's all the benefits of error types: easier, fewer SLOC, and more featureful.


Capturing and displaying multiple stacks
----------------------------------------

Each meep error can capture a stack.
This means you can have multiple stack traces reported from a single error,
if it has causes that also report stack traces:

```
Error[meep.Woop]: Wonk="Bonk";
	Caused by: Error[meep.Boop]:
		Caused by: Error[meep.Boop]:
			Stack trace:
				·> /build/path/polydawn/meep/trait_autodescriber_test.go:120: meep.TestAutodescribePlusTraceableCauseDoubleTrouble
				·> /usr/local/go/src/testing/testing.go:610: testing.tRunner
				·> /usr/local/go/src/runtime/asm_amd64.s:2086: runtime.goexit
		Stack trace:
			·> /build/path/polydawn/meep/trait_autodescriber_test.go:122: meep.TestAutodescribePlusTraceableCauseDoubleTrouble
			·> /usr/local/go/src/testing/testing.go:610: testing.tRunner
			·> /usr/local/go/src/runtime/asm_amd64.s:2086: runtime.goexit
```

As you can see, the formatting will be consistent, and each stack trace's indentation associates it with the error.


"try"-like handling and dispatch blocks
---------------------------------------

Use the `TryPlan` structure to declare error handling clearly.

[Example](https://godoc.org/github.com/polydawn/meep#example-Try):

```
meep.Try(func() {
    panic(meep.New(&meep.AllTraits{}))
}, meep.TryPlan{
    {ByType: &meep.ErrInvalidParam{},
        Handler: meep.TryHandlerMapto(&meep.ErrProgrammer{})},
    {ByVal: io.EOF,
        Handler: meep.TryHandlerDiscard},
    {CatchAny: true,
        Handler: func(error) {
            fmt.Println("caught wildcard")
        }},
})
```

There are three different ways to invoke a `TryPlan`:

- `meep.Try` (as shown above), which handles panics from the function
- [TryPlan.Handle](https://godoc.org/github.com/polydawn/meep#TryPlan.Handle), which takes the error as an argument (for use with the golang convention of returning errors)
- and [TryPlan.MustHandle](https://godoc.org/github.com/polydawn/meep#TryPlan.MustHandle), which is the same as `Handle` but will panic if there is an error and it is not explicitly handled.

`TryPlan` supports handling by **type** (which is what you'll use 99.9% of the time, since your errors with meep *are* typed)
as well as by **value**, so it works well with legacy code and existing interfaces.
You can also specify an arbitrary `func (error) (matches bool)` predicate for flexibility,
as well as a catch-all.

Note that you can specify `TryPlan` blocks separately from where you use them,
so you can declare your error handling patterns in one place and use them in *several* places.
You can also aggregate `TryPlan`s with plain ol' `append()`.

Can you still use simple built-in `switch err.(type) { ...` constructs for handling?
Yes, absolutely!  Type switches work perfectly well with meep errors.
The allure of the TryPlan is the mixture of handling type and value errors in the same block of logic.
Go with what works for you.


Availability
------------

`meep` works with any recent version of Go, going as far back as go1.2.


Performance
-----------

Quite good.

The declarative `TryPlan` structures, as shown in the example, play so nicely
with the Go compiler's escape analysis that there's effectively zero overhead to the declaration.
When evaluating a `TryPlan` with a nil error, the first branch is to return early,
so there's also effectively zero overhead in the happy/non-error path.

Capturing stack traces *is* an expensive operation.
(This is true no matter what tools or libraries you use.)
It is not recommended to mix in `TraitTraceable` in errors you return frequently or use for flow control.
Use it when you need it -- that's why it's optional.
