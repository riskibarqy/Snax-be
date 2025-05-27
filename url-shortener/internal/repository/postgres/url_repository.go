package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type urlRepository struct {
	db *pgx.Conn
}

// NewURLRepository creates a new PostgreSQL URL repository
func NewURLRepository(db *pgx.Conn) domain.URLRepository {
	return &urlRepository{
		db: db,
	}
}

func (r *urlRepository) Create(url *domain.URL) error {
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO urls (short_code, original_url, user_id, expires_at, created_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		url.ShortCode, url.OriginalURL, url.UserID, url.ExpiresAt, url.CreatedAt, url.IsActive,
	).Scan(&url.ID)

	return err
}

func (r *urlRepository) GetByShortCode(shortCode string) (*domain.URL, error) {
	url := &domain.URL{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, short_code, original_url, user_id, click_count, expires_at, created_at, is_active
		FROM urls WHERE short_code = $1`,
		shortCode,
	).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.ClickCount,
		&url.ExpiresAt, &url.CreatedAt, &url.IsActive)

	if err != nil {
		return nil, err
	}

	return url, nil
}

func (r *urlRepository) GetByUserID(userID string) ([]domain.URL, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, short_code, original_url, user_id, click_count, expires_at, created_at, is_active
		FROM urls WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []domain.URL
	for rows.Next() {
		var url domain.URL
		err := rows.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID,
			&url.ClickCount, &url.ExpiresAt, &url.CreatedAt, &url.IsActive)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (r *urlRepository) Delete(id int64, userID string) error {
	result, err := r.db.Exec(context.Background(),
		`UPDATE urls SET is_active = false WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return &domain.ErrURLNotFound{ShortCode: ""}
	}

	return nil
}

func (r *urlRepository) IncrementClickCount(id int64) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE urls SET click_count = click_count + 1 WHERE id = $1`,
		id,
	)
	return err
}
