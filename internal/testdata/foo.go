package testdata

import (
	"fmt"
	"time"
)

type Foo struct {
	noMeta     string
	NoMetaJSON string       `json:"omitempty"`
	name, Desc string       `meta:"getter"`
	size       int          `meta:"ptr;getter;setter"`
	labels     []string     `meta:"setter;getter;filter;map:time.Time"`
	stringer   fmt.Stringer `meta:"setter"`
}

func (f Foo) String() string {
	return f.name
}

type Bar struct {
	foos  []Foo              `meta:"getter;setter;map:string"`
	pairs map[string]float64 `meta:"getter;setter"`
	times []time.Time        `meta:"getter;setter;filter;map:int64"`
	baz   bool               `meta:"setter"`
}
