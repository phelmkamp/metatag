package cat

import (
	"github.com/satori/go.uuid"
)

type Cat struct {
	Uuid      uuid.UUID `meta:"getter"`
}
