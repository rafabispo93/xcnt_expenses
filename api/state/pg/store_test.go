package pg

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.azure.com/truckersb2b/edge/members.git/state"
	"dev.azure.com/truckersb2b/edge/service.git/logging"
)

func loggerCtx() context.Context {
	l := &log.Logger{Handler: discard.Default}
	return logging.WithLogger(context.TODO(), l)
}

func TestInterfaceCompliance(t *testing.T) {
	t.Parallel()
	var _ state.Store = &Store{}
}

func TestPing(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	s := &Store{
		db: dbx,
	}

	assert.NoError(t, s.Ping(context.TODO()))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_tx(t *testing.T) {
	t.Parallel()

	t.Run("Panic", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		dbx := sqlx.NewDb(db, "postgres")
		s := &Store{
			db: dbx,
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		err = s.tx(context.TODO(), func(context.Context, *sqlx.Tx) error {
			panic("some panic")
		})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		dbx := sqlx.NewDb(db, "postgres")
		s := &Store{
			db: dbx,
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		err = s.tx(context.TODO(), func(context.Context, *sqlx.Tx) error {
			return errors.New("some error")
		})
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

const (
	insertFormat = `INSERT INTO %s (%s) VALUES %%s`
	upsertFormat = `INSERT INTO %s (%s) VALUES %%s %s`
)

func int64Ptr(i int64) *int64 {
	return &i
}

func uint64Ptr(u uint64) *uint64 {
	return &u
}

func ii2dvv(ii []interface{}) []driver.Value {
	vv := make([]driver.Value, len(ii))
	for n, i := range ii {
		vv[n] = driver.Value(i)
	}
	return vv
}

func regexifyQuery(q string) string {
	replace := []struct {
		from string
		to   string
	}{
		{"(", `\(`},
		{")", `\)`},
		{"?", `\?`},
		{"$", `\$`},
		{".", `\.`},
	}

	for _, r := range replace {
		q = strings.ReplaceAll(q, r.from, r.to)
	}
	return "^" + q + "$"
}

func makePlaceholders(n int, c int) string {
	ss := make([]string, c)
	for i := 0; i < c; i++ {
		ss[i] = fmt.Sprintf("$%d", (n*c)+i+1)
	}
	return strings.Join(ss, ",")
}

func makeInserts(o int, c int) string {
	inserts := make([]string, o)
	for n := 0; n < o; n++ {
		pp := makePlaceholders(n, c)
		inserts[n] = fmt.Sprintf("(%s)", pp)
	}
	return strings.Join(inserts, ",")
}
