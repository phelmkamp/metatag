# metatag
Go metaprogramming using struct tags + generate

# Installation

`go get github.com/phelmkamp/metatag`

# Usage

1. Define struct tags

	Format is `meta:"[directive1][;directive2]"`. For example:
	```go
	type Foo struct {
		name, Desc string   `meta:"getter"`
		size       int      `meta:"ptr;getter;setter"`
		labels     []string `meta:"setter;getter;filter"`
	}
	```

2. Run command

	```bash
	metatag --path=$SRCDIR
	```

	Better yet, add the following comment to a file at the root of your source tree (e.g. main.go)
	and run `go generate` as part of your build process.

	```go
	//go:generate metatag
	```

3. Enjoy!

	A *_meta.go file is generated for each *.go file that has meta tags.
	You can review/modify the generated code, write corresponding tests, etc!
	Just be aware that any changes will be overwritten the next time the tool runs.

# Directives

`getter`

Generates a getter. Method name is the name of the field.
`Get` is prepended to the name if and only if the field is already exported.
Uses value receiver by default.

`setter`

Generates a setter. Method name is the name of the field prepended with `Set`.
Always uses pointer receiver.

`filter` (slice only)

Generates a method that returns a copy of the slice, omitting elements that are rejected by the given function.
Method name is `Filter`, followed by the name of the field unless `omitfield` is specified.
Uses value receiver by default.

`map:$Type` (slice only)

Generates a method that returns a copy of the slice, mapping elements to the specified type using the given function.
Method name is of the form `MapFieldTo$Type`, or just `MapTo$Type` if `omitfield` is specified.
Uses value receiver by default.

`stringer`

Includes the field in the result of the generated `String` method. Uses value receiver by default.

`new`

Includes the field as an argument to the generated `NewType` method.

`ptr`

Specifies that a pointer receiver be used for all subsequent directives.

# FAQ

1. Why generate getters and setters?

	Getters/setters are sometimes necessary to adhere to a "data contract" since [Go interfaces only match methods, not fields](https://github.com/golang/go/issues/23796).
	Getters/setters are great candidates for code generation because they are true boilerplate where names and types directly correspond to a particular field.

2. Why code generation instead of reflection?

	Code generation provides compile-time type safety which is a critical feature of Go and languages like it.
	Plus, generation produces easy-to-understand code that you can review and modify as you see fit!

3. Why struct tags?

	Struct tags are well suited for this task because they are designed to provide auxilliary information to tools/packages in a concise and unobtrusive way.
	Also, generating methods for a struct gives us a nice "namespace" with a low probability of collisions (as opposed to package-level functions).