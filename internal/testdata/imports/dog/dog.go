package dog

import (
	"github.com/satori/go.uuid"
)

type Dog struct {
	Uuid      uuid.UUID `meta:"equal"`
}
