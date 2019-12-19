package person

import "time"

type Person struct {
	Name      string `meta:"stringer"`
	Birthdate time.Time
}

type Persons struct {
	Persons []Person `meta:"filter,omitfield;map:int,omitfield"`
}
