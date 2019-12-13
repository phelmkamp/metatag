package testdata

import (
	"fmt"
	"time"
)

type Foo struct {
	noMeta     string
	NoMetaJSON string       `json:"omitempty"`
	name, Desc string       `meta:"getter;stringer"`
	size       int          `meta:"ptr;getter;setter;stringer"`
	labels     []string     `meta:"setter;getter;filter;map:time.Time"`
	stringer   fmt.Stringer `meta:"setter"`
}

type Bar struct {
	foos  []Foo              `meta:"getter;setter;map:string"`
	pairs map[string]float64 `meta:"getter;setter"`
	times []time.Time        `meta:"getter;setter;filter;map:int64"`
	baz   bool               `meta:"setter"`
}
