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

	Better yet, add the following comment to your source code (e.g. main.go) and run `go generate` as part of your build process.

	```go
	//go:generate metatag
	```

3. Enjoy!

	A *_meta.go file is generated for each *.go file that has meta tags. You can review/modify the generated code, write corresponding tests, etc.

# Directives
`getter`

Generates a getter. Method name is the name of the field. `Get` is prepended to the name if and only if the field is already exported. Uses value receiver by default.

`setter`

Generates a setter. Method name is the name of the field prepended with `Set`. Always uses pointer receiver.

`filter` (slice only)

Generates a method that returns a copy of the slice, omitting elements that are rejected by the given function. Method name is the name of the field prepended with `Filter`. Uses value receiver by default.

`map:$type` (slice only)

Generates a method that returns a copy of the slice, mapping elements to the specified type using the given function. Method name is of the form `MapFieldToType`. Uses value receiver by default.

`stringer`

stringer stringer stringer stringer

`new`

new new new new

`ptr`

Specifies that a pointer receiver be used for all subsequent directives.
