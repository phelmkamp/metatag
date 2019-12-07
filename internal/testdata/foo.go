package testdata

type Foo struct {
	name, Desc string   `meta:"getter"`
	size       int      `meta:"ptr;getter;setter"`
	labels     []string `meta:"setter;getter;find;filter"`
}

type Bar struct {
	foos  []Foo              `meta:"getter;setter;filter;find"`
	pairs map[string]float64 `meta:"getter;setter"`
}
