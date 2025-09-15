package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	// url *model.URL
	db *sql.DB
}

func NewPostgresRepository(ctx context.Context, connectionString string) *PostgresRepository {
	var err error
	pg := &PostgresRepository{}
	
	inctx, cansel := context.WithTimeout(ctx, 10 * time.Second)
	defer cansel()

	pg.db, err = sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := pg.Ping(inctx); err != nil {
		log.Fatal(err.Error())	
	}
	pg.LoadingRepository(inctx)
	return pg
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
	insertSQL := `INSERT INTO a_url_short (uuid, short_url, original_url) VALUES ($1, $2, $3)`
	if _, err := p.db.Exec(insertSQL, url.UUID, url.ShortURL, url.OriginalURL); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			errSelect := p.db.QueryRowContext(ctx, 
				"select uuid, short_url from a_url_short WHERE original_url = $1",
				url.OriginalURL,
			).Scan(&url.UUID, &url.ShortURL)
			if errSelect != nil {
				return nil, errSelect
			}
			return url, model.ErrURLConflict
		}
		return nil, err
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
	createTable := `CREATE TABLE IF NOT EXISTS a_url_short (id SERIAL PRIMARY KEY,uuid text not null,short_url text not null,original_url text UNIQUE not null);CREATE index IF NOT EXISTS idx_short_url_uuid ON a_url_short(uuid);`
	if _ , err := p.db.ExecContext(ctx, createTable); err != nil {
		return err
	}
	return nil
}
