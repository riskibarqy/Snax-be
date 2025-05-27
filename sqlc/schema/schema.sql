-- This file is auto-generated from migrations. DO NOT EDIT DIRECTLY.
-- Last updated: Wed, May 28, 2025 12:15:56 AM

-- Including migration: 000001_create_urls_table.up.sql

-- Create URLs table
CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    click_count BIGINT NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);
CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at); 

-- Including migration: 000002_create_analytics_table.up.sql

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

-- Including migration: 000003_create_tags_tables.up.sql

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

-- Including migration: 000004_create_custom_domains.up.sql

-- Create custom domains table
CREATE TABLE IF NOT EXISTS custom_domains (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL UNIQUE,
    user_id VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create index
CREATE INDEX IF NOT EXISTS idx_custom_domains_user ON custom_domains(user_id); 

