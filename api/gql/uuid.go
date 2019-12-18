package gql

import (
	"io"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// UUID is a custom scalar wrapping a uuid.UUID.
type UUID uuid.UUID

// UnmarshalGQL takes a string and converts it to a UUID.
func (u *UUID) UnmarshalGQL(v interface{}) error {
	if s, ok := v.(string); ok {
		id, err := uuid.FromString(s)
		if err != nil {
			return errors.Wrapf(err, "could not parse input into UUID: %s", s)
		}

		*u = UUID(id)
		return nil
	}

	return errors.Errorf("could not parse input into UUID: %[1]v (%[1]T)", v)
}

// MarshalGQL writes a UUID to the writer.
func (u UUID) MarshalGQL(w io.Writer) {
	w.Write([]byte(`"` + uuid.UUID(u).String() + `"`))
}
