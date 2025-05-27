package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riskibarqy/Snax-be/url-shortener/internal/domain"
)

type tagRepository struct {
	db *pgx.Conn
}

// NewTagRepository creates a new PostgreSQL tag repository
func NewTagRepository(db *pgx.Conn) domain.TagRepository {
	return &tagRepository{
		db: db,
	}
}

func (r *tagRepository) Create(tag *domain.Tag) error {
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO tags (name) VALUES ($1) RETURNING id`,
		tag.Name,
	).Scan(&tag.ID)

	return err
}

func (r *tagRepository) GetByID(id int64) (*domain.Tag, error) {
	tag := &domain.Tag{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, name FROM tags WHERE id = $1`,
		id,
	).Scan(&tag.ID, &tag.Name)

	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *tagRepository) GetByName(name string) (*domain.Tag, error) {
	tag := &domain.Tag{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, name FROM tags WHERE name = $1`,
		name,
	).Scan(&tag.ID, &tag.Name)

	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *tagRepository) GetByURLID(urlID int64) ([]domain.Tag, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT t.id, t.name FROM tags t
		JOIN url_tags ut ON ut.tag_id = t.id
		WHERE ut.url_id = $1`,
		urlID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *tagRepository) AddURLTag(urlID, tagID int64) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO url_tags (url_id, tag_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING`,
		urlID, tagID,
	)
	return err
}

func (r *tagRepository) RemoveURLTag(urlID, tagID int64) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM url_tags WHERE url_id = $1 AND tag_id = $2`,
		urlID, tagID,
	)
	return err
}

func (r *tagRepository) AddTagToURL(ctx context.Context, urlID int64, tagName string) error {
	// First get or create the tag
	tag, err := r.GetByName(tagName)
	if err != nil {
		// Create new tag if it doesn't exist
		tag = &domain.Tag{Name: tagName}
		if err := r.Create(tag); err != nil {
			return err
		}
	}

	// Add the tag to URL
	return r.AddURLTag(urlID, tag.ID)
}

func (r *tagRepository) RemoveTagFromURL(ctx context.Context, urlID int64, tagName string) error {
	// Get the tag
	tag, err := r.GetByName(tagName)
	if err != nil {
		return err
	}

	// Remove the tag from URL
	return r.RemoveURLTag(urlID, tag.ID)
}

func (r *tagRepository) GetURLTags(ctx context.Context, urlID int64) ([]domain.Tag, error) {
	return r.GetByURLID(urlID)
}
