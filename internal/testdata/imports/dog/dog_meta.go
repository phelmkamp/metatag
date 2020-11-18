// GENERATED BY metatag, DO NOT EDIT
// (or edit away - I'm a comment, not a cop)

package dog

import (
	"github.com/satori/go.uuid"
)

// Equal answers whether v is equivalent to d.
// Always returns false if v is not a Dog.
func (d Dog) Equal(v interface{}) bool {
	d2, ok := v.(Dog)
	if !ok {
		return false
	}
	if d.Uuid != d2.Uuid {
		return false
	}
	return true
}
