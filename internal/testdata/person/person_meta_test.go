package person

import (
	"fmt"
	"time"
)

func Example() {
	persons := Persons{[]Person{
		{
			Name:      "Ann",
			Birthdate: time.Date(1983, time.February, 12, 0, 0, 0, 0, time.Local),
		},
		{
			Name:      "Bob",
			Birthdate: time.Date(2007, time.August, 8, 0, 0, 0, 0, time.Local),
		},
		{
			Name:      "Charlie",
			Birthdate: time.Date(1994, time.June, 28, 0, 0, 0, 0, time.Local),
		},
	}}

	var name string
	hasName := func(p Person) bool { return p.Name == name }

	name = "David"
	if found := persons.FilterPersons(hasName, 1); len(found) > 0 {
		// contains David
		fmt.Println(found[0])
	}

	name = "Bob"
	if found := persons.FilterPersons(hasName, 1); len(persons.FilterPersons(hasName, 1)) > 0 {
		// contains Bob
		fmt.Println(found[0])
	}

	// exludes Bob
	found := persons.FilterPersons(func(p Person) bool { return p.Name != name }, -1)
	fmt.Println(found)

	// map to ages
	ages := persons.MapPersonsToInt(func(p Person) int { return time.Now().Year() - p.Birthdate.Year() })
	fmt.Println(ages)

	// Output: Bob
	// [Ann Charlie]
	// [36 12 25]
}
