-- Create analytics table
CREATE TABLE IF NOT EXISTS analytics (
    id SERIAL PRIMARY KEY,
    url_id INTEGER REFERENCES urls(id),
    visitor_ip VARCHAR(45),
    user_agent TEXT,
    referer TEXT,
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    country_code VARCHAR(2),
    device_type VARCHAR(20)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_url_id ON analytics(url_id);
CREATE INDEX IF NOT EXISTS idx_timestamp ON analytics(timestamp); 