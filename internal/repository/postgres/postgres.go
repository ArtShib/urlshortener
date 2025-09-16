package postgres

import (
	"context"
	"database/sql"

	"github.com/ArtShib/urlshortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(ctx context.Context, connectionString string) (*PostgresRepository, error) {
	var err error
	pg := &PostgresRepository{}

	pg.db, err = sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	if err := pg.Ping(ctx); err != nil {
		return nil, err	
	}
	pg.LoadingRepository(ctx)
	return pg, nil
}

func (p PostgresRepository) Ping(ctx context.Context) error {
	if err := p.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
func (p *PostgresRepository) Close() error {
	if err := p.db.Close(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) Save(ctx context.Context, url *model.URL) (*model.URL, error) {
	var isConflict bool
	insertSQL :=  `WITH inserted AS (
						INSERT INTO a_url_short (uuid, short_url, original_url)
						VALUES ($1, $2, $3)
						ON CONFLICT (original_url) DO NOTHING
						RETURNING *
					)
					select uuid, short_url, false as is_conflict FROM inserted
					UNION
					SELECT uuid, short_url, true as is_conflict FROM a_url_short 
					WHERE original_url = $3 AND NOT EXISTS (SELECT 1 FROM inserted)`
	err := p.db.QueryRowContext(ctx, insertSQL, url.UUID, url.ShortURL, url.OriginalURL).
			Scan(&url.UUID, &url.ShortURL, &isConflict)

	if err != nil {
		return nil, err
	}
	if isConflict {
		return url, model.ErrURLConflict 
	}
	return url, nil
}

func (p *PostgresRepository) Get(ctx context.Context, uuid string) (*model.URL, error) {
	query := `select uuid, short_url, original_url from a_url_short where uuid = $1 LIMIT 1`
	row := p.db.QueryRowContext(ctx, query, uuid)
	var url model.URL
	if err := row.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL); err != nil {
		return nil, err
	}
	return &url, nil
}

func (p *PostgresRepository) LoadingRepository(ctx context.Context) error {
	createTable := `CREATE TABLE IF NOT EXISTS a_url_short (
						id SERIAL PRIMARY KEY,
						uuid text not null,
						short_url text not null,
						original_url text UNIQUE not null);
					CREATE index IF NOT EXISTS idx_short_url_uuid ON a_url_short(uuid);`
	if _ , err := p.db.ExecContext(ctx, createTable); err != nil {
		return err
	}
	return nil
}
