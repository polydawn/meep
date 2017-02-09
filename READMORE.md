
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
