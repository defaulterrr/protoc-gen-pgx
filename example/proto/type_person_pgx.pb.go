package pb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonsRepository struct {
	pool *pgxpool.Pool
}

func NewPersonsRepository(pool *pgxpool.Pool) *PersonsRepository {
	return &PersonsRepository{
		pool: pool,
	}
}

type Persons struct {
	name string
	id   int64
}

const insertQuery = `
INSERT INTO persons 
(name,id)
VALUES 
(@name, @id)
`

func (repo *PersonsRepository) InsertPersons(ctx context.Context, persons Persons) error {
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	args := pgx.NamedArgs{
		"name": persons.name,
		"id":   persons.id,
	}

	if _, err := conn.Exec(ctx, insertQuery, args); err != nil {
		return fmt.Errorf("failed to insert data into persons table: %w", err)
	}

	return nil
}
