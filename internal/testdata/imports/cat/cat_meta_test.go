package cat

import (
	"github.com/satori/go.uuid"
)

func Example() {
	cat := Cat{
		Uuid:      uuid.NewV1(),
	}

	cat.GetUuid()
}

