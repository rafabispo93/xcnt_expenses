package gql

import (
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	t.Parallel()

	t.Run("UnmarshalGQL", func(t *testing.T) {
		t.Parallel()

		id := uuid.NewV4()
		var u UUID
		require.NoError(t, u.UnmarshalGQL(id.String()))

		assert.Equal(t, id, uuid.UUID(u))
	})

	t.Run("MarshalGQL", func(t *testing.T) {
		t.Parallel()

		id := uuid.NewV4()
		u := UUID(id)

		var sb strings.Builder
		u.MarshalGQL(&sb)

		assert.Equal(t, `"` + id.String() + `"`, sb.String())
	})
}

