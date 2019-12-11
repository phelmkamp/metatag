package testdata

import "time"

import "fmt"

type Foo struct {
	noMeta     string
	NoMetaJSON string       `json:"omitempty"`
	name, Desc string       `meta:"getter"`
	size       int          `meta:"ptr;getter;setter"`
	labels     []string     `meta:"setter;getter;find;filter"`
	stringer   fmt.Stringer `meta:"setter"`
}

type Bar struct {
	foos  []Foo              `meta:"getter;setter;filter;find"`
	pairs map[string]float64 `meta:"getter;setter"`
	times []time.Time        `meta:"getter;setter;filter;find"`
	baz   bool               `meta:"setter"`
}
