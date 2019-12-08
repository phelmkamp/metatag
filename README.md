# metatag
Go metaprogramming using struct tags + generate (no reflection!)

# Installation
`go get github.com/phelmkamp/metatag`

# Usage
1. Define struct tags

	Format is `meta:"[directive1][;directive2]"`. For example:
	```go
	type Foo struct {
		name, Desc string   `meta:"getter"`
		size       int      `meta:"ptr;getter;setter"`
		labels     []string `meta:"setter;getter;find;filter"`
	}
	```

2. Run command
	```bash
	metatag --path=$SRCDIR
	```

	For best results, add the following comment to your source code (e.g. main.go) and run `go generate`.

	```go
	//go:generate metatag
	```

# Directives
`getter`

Generates a getter. Method name is the name of the field. `Get` is prepended to the name if and only if the field is already exported. Uses value receiver by default.

`setter`

Generates a setter. Method name is the name of the field prepended with `Set`. Always uses pointer receiver.

`filter` (slice only)

Generates a method that returns a copy of the slice, omitting elements that are rejected by the given function. Method name is the name of the field prepended with `Filter`. Uses value receiver by default.

`find` (slice only)

Generates a method that returns the index of the first element that matches the argument (using `reflect.DeepEqual`). Method name is the name of the field minus any plural `s`, prepended with `Find`. Uses value receiver by default.

`ptr`

Specifies that a pointer receiver be used for all subsequent directives.
