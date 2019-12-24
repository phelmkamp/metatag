package person

import "time"

type Person struct {
	Name      string `meta:"stringer"`
	Birthdate time.Time
}

type Persons struct {
	result []Person `meta:"wrapper;new;filter;mapper,int;sort,stringer;getter"`
}
