package domain

// Tag represents a URL tag
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TagRepository defines the interface for tag data operations
type TagRepository interface {
	Create(tag *Tag) error
	GetByID(id int64) (*Tag, error)
	GetByName(name string) (*Tag, error)
	GetByURLID(urlID int64) ([]Tag, error)
	AddURLTag(urlID, tagID int64) error
	RemoveURLTag(urlID, tagID int64) error
}

// TagService defines the interface for tag business logic
type TagService interface {
	CreateTag(name string) (*Tag, error)
	GetTag(id int64) (*Tag, error)
	GetTagByName(name string) (*Tag, error)
	GetURLTags(urlID int64) ([]Tag, error)
	AddTagToURL(urlID int64, tagName string) error
	RemoveTagFromURL(urlID int64, tagName string) error
}
