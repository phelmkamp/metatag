package dog

import (
	"github.com/satori/go.uuid"
)

func Example() {
	dog := Dog{
		Uuid:      uuid.NewV1(),
	}

	dog.Equal(dog)
}

