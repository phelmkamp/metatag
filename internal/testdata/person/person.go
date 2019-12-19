package person

import "time"

type Person struct {
	Name      string `meta:"stringer"`
	Birthdate time.Time
}

type Persons struct {
	Ps []Person `meta:"filter,omitfield;map:int,omitfield;sort,stringer"`
}
