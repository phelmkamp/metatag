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
	labels     []string     `meta:"new;setter;getter;filter;mapper,time.Time"`
	stringer   fmt.Stringer `meta:"setter"`
}

type Bar struct {
	name  string             `meta:"stringer;equal"`
	foos  []Foo              `meta:"getter;setter;mapper,string"`
	pairs map[string]float64 `meta:"getter;setter"`
	times []time.Time        `meta:"getter;setter;filter;mapper,int64;equal,reflect"`
	baz   bool               `meta:"setter"`
}
