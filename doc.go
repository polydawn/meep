/*
	More Expressive Error Patterns.

	Embed `meep` types in your errors structures to effortlessly get fancy behaviors.
	Stacks, automatic message generation from your fields, common handling:
	declaring useful error types is now *much* easier.
*/
package meep

/*
	Implementation pattern notes:

	All the public stuff is made available by embedding it into a struct
	of the user-library's definition.
	These are all the in `mixins.go` file.

	Many of the mixins implement an interface which is private to this package
	and simply returns a mutable reference to themselves.
	This is so we can do type switching against that interface (instead
	needing reflection), and then use it to reach in and initialize the value.
*/
