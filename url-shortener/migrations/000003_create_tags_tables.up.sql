-- Create tags table for URL categorization
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- Create URL-tag relationship table
CREATE TABLE IF NOT EXISTS url_tags (
    url_id INTEGER REFERENCES urls(id),
    tag_id INTEGER REFERENCES tags(id),
    PRIMARY KEY (url_id, tag_id)
); 