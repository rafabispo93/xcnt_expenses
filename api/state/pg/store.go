package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const parameterLimit = 65535

// Store provides data interaction methods over a Postgres DB.
type Store struct {
	db     *sqlx.DB
	makeID func() uuid.UUID
}

// NewStore returns a configured *Store.
func NewStore(cxn string) (*Store, error) {
	db, err := sqlx.Connect("postgres", cxn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to DB")
	}

	return &Store{
		db:     db,
		makeID: uuid.NewV4,
	}, nil
}

// Ping pings the DB.
func (s *Store) Ping(ctx context.Context) error {
	return errors.Wrap(s.db.PingContext(ctx), "failed to ping DB")
}

func (s *Store) tx(ctx context.Context, fn func(context.Context, *sqlx.Tx) error) (err error) {
	var tx *sqlx.Tx
	if tx, err = s.db.BeginTxx(ctx, nil); err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("store encountered unhandled panic during transaction: %v", r)
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	err = fn(ctx, tx)
	return err
}

func makeUpsert(fields []string, columns []string) string {
	u := fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %%s", strings.Join(fields, ","))

	sets := make([]string, len(columns))
	for i, c := range columns {
		sets[i] = fmt.Sprintf("%[1]s = excluded.%[1]s", c)
	}

	return fmt.Sprintf(u, strings.Join(sets, ","))
}

func makeUpsertWithUpdated(fields []string, columns []string) string {
	return fmt.Sprintf("%s,updated = now()", makeUpsert(fields, columns))
}
