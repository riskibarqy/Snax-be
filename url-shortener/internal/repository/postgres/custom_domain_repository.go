package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	internalDomain "github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type customDomainRepository struct {
	db *pgx.Conn
}

// NewCustomDomainRepository creates a new PostgreSQL custom domain repository
func NewCustomDomainRepository(db *pgx.Conn) internalDomain.CustomDomainRepository {
	return &customDomainRepository{
		db: db,
	}
}

func (r *customDomainRepository) Create(domain *internalDomain.CustomDomain) error {
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO custom_domains (domain, user_id, verified)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		domain.Domain, domain.UserID, domain.Verified,
	).Scan(&domain.ID, &domain.CreatedAt)

	return err
}

func (r *customDomainRepository) GetByDomain(domain string) (*internalDomain.CustomDomain, error) {
	d := &internalDomain.CustomDomain{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, domain, user_id, verified, created_at
		FROM custom_domains WHERE domain = $1`,
		domain,
	).Scan(&d.ID, &d.Domain, &d.UserID, &d.Verified, &d.CreatedAt)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func (r *customDomainRepository) GetByUserID(userID string) ([]internalDomain.CustomDomain, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, domain, user_id, verified, created_at
		FROM custom_domains WHERE user_id = $1
		ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []internalDomain.CustomDomain
	for rows.Next() {
		var d internalDomain.CustomDomain
		err := rows.Scan(&d.ID, &d.Domain, &d.UserID, &d.Verified, &d.CreatedAt)
		if err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}

	return domains, nil
}

func (r *customDomainRepository) UpdateVerificationStatus(id int64, verified bool) error {
	result, err := r.db.Exec(context.Background(),
		`UPDATE custom_domains SET verified = $1 WHERE id = $2`,
		verified, id,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return &internalDomain.ErrDomainNotFound{Domain: ""}
	}

	return nil
}

func (r *customDomainRepository) Delete(id int64, userID string) error {
	result, err := r.db.Exec(context.Background(),
		`DELETE FROM custom_domains WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return &internalDomain.ErrDomainNotFound{Domain: ""}
	}

	return nil
}
