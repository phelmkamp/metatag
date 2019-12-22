package foobar

import (
	"fmt"
	"time"
)

type Foo struct {
	noMeta     string
	NoMetaJSON string       `json:"omitempty"`
	name, Desc string       `meta:"new;getter;stringer"`
	size       int          `meta:"stringer;ptr;getter;setter"`
	labels     []string     `meta:"new;setter;getter;filter;map:time.Time"`
	stringer   fmt.Stringer `meta:"setter"`
}

type Bar struct {
	name  string             `meta:"stringer;equal"`
	foos  []Foo              `meta:"getter;setter;map:string"`
	pairs map[string]float64 `meta:"getter;setter"`
	times []time.Time        `meta:"getter;setter;filter;map:int64;equal,reflect"`
	baz   bool               `meta:"setter"`
}
