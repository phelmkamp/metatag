package person

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Example() {
	ps := []Person{
		{
			Name:      "Charlie",
			Birthdate: time.Date(1994, time.June, 28, 0, 0, 0, 0, time.Local),
		},
		{
			Name:      "Bob",
			Birthdate: time.Date(2007, time.August, 8, 0, 0, 0, 0, time.Local),
		},
		{
			Name:      "Ann",
			Birthdate: time.Date(1983, time.February, 12, 0, 0, 0, 0, time.Local),
		},
	}

	var name string
	hasName := func(p Person) bool { return p.Name == name }

	name = "David"
	if found := NewPersons(ps).FilterN(hasName, 1).Result(); len(found) > 0 {
		// contains David
		fmt.Println(found[0])
	}

	name = "Bob"
	if found := NewPersons(ps).FilterN(hasName, 1).Result(); len(found) > 0 {
		// contains Bob
		fmt.Println(found[0])
	}

	ages := NewPersons(ps).
		// exclude Bob
		Filter(func(p Person) bool {
			return p.Name != "Bob"
		}).
		// sort by name
		Sort(func(vi, vj Person) bool {
			return vi.Name < vj.Name
		}).
		// map to ages
		MapToInt(func(p Person) int {
			return time.Now().Year() - p.Birthdate.Year()
		})
	fmt.Println(ages)

	// Output: Bob
	// [36 25]
}

func TestPersons_Sort(t *testing.T) {
	tests := []struct {
		name string
		p    Persons
		want []Person
	}{
		{
			name: "abc",
			p: NewPersons(
				[]Person{
					{Name: "b"},
					{Name: "c"},
					{Name: "a"},
				},
			),
			want: []Person{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Sort(func(vi, vj Person) bool {
				return vi.Name < vj.Name
			})
			if !reflect.DeepEqual(tt.p.Result(), tt.want) {
				t.Errorf("got = %v, want %v", tt.p.Result(), tt.want)
			}
		})
	}
}
