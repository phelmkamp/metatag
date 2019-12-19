package person

import (
	"fmt"
	"reflect"
	"testing"
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
	if found := persons.Filter(hasName, 1); len(found) > 0 {
		// contains David
		fmt.Println(found[0])
	}

	name = "Bob"
	if found := persons.Filter(hasName, 1); len(found) > 0 {
		// contains Bob
		fmt.Println(found[0])
	}

	// exludes Bob
	found := persons.Filter(func(p Person) bool { return p.Name != name }, -1)
	fmt.Println(found)

	// map to ages
	ages := persons.MapToInt(func(p Person) int { return time.Now().Year() - p.Birthdate.Year() })
	fmt.Println(ages)

	// Output: Bob
	// [Ann Charlie]
	// [36 12 25]
}

func TestPersons_Sort(t *testing.T) {
	tests := []struct {
		name string
		p    Persons
		want []Person
	}{
		{
			name: "abc",
			p: Persons{[]Person{
				{Name: "b"},
				{Name: "c"},
				{Name: "a"},
			}},
			want: []Person{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Sort()
			if !reflect.DeepEqual(tt.p.Ps, tt.want) {
				t.Errorf("got = %v, want %v", tt.p.Ps, tt.want)
			}
		})
	}
}
