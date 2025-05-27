package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type analyticsRepository struct {
	db *pgx.Conn
}

// NewAnalyticsRepository creates a new PostgreSQL analytics repository
func NewAnalyticsRepository(db *pgx.Conn) domain.AnalyticsRepository {
	return &analyticsRepository{
		db: db,
	}
}

func (r *analyticsRepository) Create(analytics *domain.Analytics) error {
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO analytics (url_id, visitor_ip, user_agent, referer, country_code, device_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, timestamp`,
		analytics.URLID, analytics.VisitorIP, analytics.UserAgent, analytics.Referer,
		analytics.CountryCode, analytics.DeviceType,
	).Scan(&analytics.ID, &analytics.Timestamp)

	return err
}

func (r *analyticsRepository) GetByURLID(urlID int64) ([]domain.Analytics, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, url_id, visitor_ip, user_agent, referer, timestamp, country_code, device_type
		FROM analytics WHERE url_id = $1 ORDER BY timestamp DESC`,
		urlID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []domain.Analytics
	for rows.Next() {
		var a domain.Analytics
		err := rows.Scan(
			&a.ID, &a.URLID, &a.VisitorIP, &a.UserAgent, &a.Referer,
			&a.Timestamp, &a.CountryCode, &a.DeviceType,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}

	return analytics, nil
}
