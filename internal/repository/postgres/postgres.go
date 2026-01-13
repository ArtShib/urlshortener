package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ArtShib/urlshortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// RepositoryPostgres структура для работы с БД
type RepositoryPostgres struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewPostgresRepository конструтор для RepositoryPostgres
func NewPostgresRepository(ctx context.Context, connectionString string, log *slog.Logger) (*RepositoryPostgres, error) {
	const op = "postgres.NewPostgresRepository"
	logger := log.With(
		slog.String("op", op),
	)

	var err error
	pg := &RepositoryPostgres{
		logger: log,
	}

	pg.db, err = sql.Open("pgx", connectionString)
	if err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := pg.Ping(ctx); err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := pg.LoadingRepository(ctx); err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return pg, nil
}

// Ping метод для проверки доступности БД
func (p *RepositoryPostgres) Ping(ctx context.Context) error {
	const op = "postgres.Ping"
	logger := p.logger.With(
		slog.String("op", op),
	)
	if err := p.db.PingContext(ctx); err != nil {
		logger.Error(op, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Close метод для закрытия соединения с БД
func (p *RepositoryPostgres) Close() error {
	const op = "postgres.Close"
	logger := p.logger.With(
		slog.String("op", op),
	)
	if err := p.db.Close(); err != nil {
		logger.Error(op, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Save метод для сохрания сокращенного url
func (p *RepositoryPostgres) Save(ctx context.Context, url *model.URL) (*model.URL, error) {
	const op = "postgres.Close"
	logger := p.logger.With(
		slog.String("op", op),
	)
	var isConflict bool
	stmt, err := p.db.Prepare(`WITH inserted AS (
						INSERT INTO a_url_short (uuid, short_url, original_url, user_id)
						VALUES ($1, $2, $3, $4)
						ON CONFLICT (original_url) DO NOTHING
						RETURNING *
					)
					select uuid, short_url, false as is_conflict FROM inserted
					UNION
					SELECT uuid, short_url, true as is_conflict FROM a_url_short 
					WHERE original_url = $3 AND NOT EXISTS (SELECT 1 FROM inserted)`)

	if err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := stmt.QueryRowContext(ctx, url.UUID, url.ShortURL, url.OriginalURL, url.UserID).Scan(&url.UUID, &url.ShortURL, &isConflict); err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if isConflict {
		return url, model.ErrURLConflict
	}
	return url, nil
}

// Get метод получения оригинального url
func (p *RepositoryPostgres) Get(ctx context.Context, uuid string) (*model.URL, error) {
	const op = "postgres.Get"
	logger := p.logger.With(
		slog.String("op", op),
	)
	stmt, err := p.db.Prepare(`select uuid, short_url, original_url, user_id, is_deleted from a_url_short where uuid = $1 LIMIT 1`)
	if err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, uuid)
	var url model.URL
	if err := row.Scan(&url.UUID, &url.ShortURL, &url.OriginalURL, &url.UserID, &url.DeletedFlag); err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &url, nil
}

// LoadingRepository метод подготовки БД
func (p *RepositoryPostgres) LoadingRepository(ctx context.Context) error {
	createTable := `CREATE TABLE IF NOT EXISTS a_url_short (
						id SERIAL PRIMARY KEY,
						uuid text not null,
						short_url text not null,
						original_url text UNIQUE not null,
						user_id text default null,
                        is_deleted boolean default false);
					CREATE index IF NOT EXISTS idx_short_url_uuid ON a_url_short(uuid);`
	if _, err := p.db.ExecContext(ctx, createTable); err != nil {
		return err
	}
	return nil
}

// GetBatch метод получения оригинального url по id пользователя
func (p *RepositoryPostgres) GetBatch(ctx context.Context, userID string) (model.URLUserBatch, error) {
	const op = "postgres.GetBatch"
	logger := p.logger.With(
		slog.String("op", op),
	)
	stmt, err := p.db.Prepare(`select short_url, original_url from a_url_short where user_id = $1`)
	if err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if rows == nil {
		return model.URLUserBatch{}, nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error(op, "error", err)
		}
	}()

	var urls model.URLUserBatch
	for rows.Next() {
		var url model.URLUser
		if err := rows.Scan(&url.ShortURL, &url.OriginalURL); err != nil {
			logger.Error(op, "error", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return model.URLUserBatch{}, fmt.Errorf("%s: %w", op, err)
	}
	return urls, nil
}

// DeleteBatch метод установки признака удаления url
func (p *RepositoryPostgres) DeleteBatch(ctx context.Context, deleteRequest model.URLUserRequestArray) error {
	const op = "postgres.DeleteBatch"
	logger := p.logger.With(
		slog.String("op", op),
	)

	values := make([]string, len(deleteRequest))
	args := make([]interface{}, 0, len(deleteRequest)*2)
	for i, req := range deleteRequest {
		pos1, pos2 := len(args)+1, len(args)+2
		values[i] = fmt.Sprintf("($%d, $%d)", pos1, pos2)
		args = append(args, req.UUID, req.UserID)
	}

	stmt, err := p.db.Prepare(fmt.Sprintf(`
        UPDATE a_url_short 
        SET is_deleted = true
        FROM (VALUES %s) AS targets(uuid, user_id)
        WHERE a_url_short.uuid = targets.uuid 
          AND a_url_short.user_id = targets.user_id
          AND a_url_short.is_deleted = false`,
		strings.Join(values, ", ")))
	if err != nil {
		logger.Error(op, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		logger.Error(op, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
