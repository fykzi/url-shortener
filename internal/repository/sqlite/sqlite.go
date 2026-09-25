package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fykzi/url-shortener/internal/model"
)

type LinkRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewLinkRepository(db *sql.DB, log *slog.Logger) *LinkRepository {
	return &LinkRepository{
		db:  db,
		log: log,
	}
}

func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	const query = `
		INSERT INTO links (short_code, original_url, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, link.ShortCode, link.OriginalURL, link.CreatedAt, link.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}

func (r *LinkRepository) GetByCode(ctx context.Context, code string) (*model.Link, error) {
	const query = `
		SELECT id, short_code, original_url, clicks, created_at, updated_at
		FROM links
		WHERE short_code = ?
	`

	var link model.Link

	err := r.db.QueryRowContext(ctx, query, code).
		Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.Clicks, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get link by code: %w", err)
	}

	return &link, nil
}

func (r *LinkRepository) Delete(ctx context.Context, code string) error {
	const query = `DELETE FROM links WHERE short_code = ?`

	result, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *LinkRepository) List(ctx context.Context) ([]*model.Link, error) {
	const query = `
		SELECT id, short_code, original_url, clicks, created_at, updated_at
		FROM links
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list links: %w", err)
	}
	defer rows.Close()

	var links []*model.Link

	for rows.Next() {
		var link model.Link

		if err := rows.Scan(
			&link.ID,
			&link.ShortCode,
			&link.OriginalURL,
			&link.Clicks,
			&link.CreatedAt,
			&link.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan link: %w", err)
		}
		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return links, nil
}

func (r *LinkRepository) Redirect(ctx context.Context, code string) (string, error) {
	const query = `
		UPDATE links
		SET clicks = clicks + 1, updated_at = CURRENT_TIMESTAMP
		WHERE short_code = ?
		RETURNING original_url
	`

	var originalURL string

	if err := r.db.QueryRowContext(ctx, query, code).Scan(&originalURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("redirect: %w", err)
	}

	return originalURL, nil
}
