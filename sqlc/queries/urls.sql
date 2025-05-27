-- name: CreateURL :one
INSERT INTO urls (
    short_code,
    original_url,
    user_id,
    expires_at
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetURLByShortCode :one
SELECT * FROM urls
WHERE short_code = $1 AND is_active = true;

-- name: GetURLsByUserID :many
SELECT * FROM urls
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: DeactivateURL :exec
UPDATE urls
SET is_active = false
WHERE id = $1 AND user_id = $2;

-- name: IncrementClickCount :one
UPDATE urls
SET click_count = click_count + 1
WHERE id = $1
RETURNING click_count;

-- name: CreateAnalytics :one
INSERT INTO analytics (
    url_id,
    visitor_ip,
    user_agent,
    referer,
    country_code,
    device_type
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAnalyticsByURLID :many
SELECT * FROM analytics
WHERE url_id = $1
ORDER BY timestamp DESC;

-- name: AddURLTag :exec
INSERT INTO url_tags (url_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: GetURLTags :many
SELECT t.* FROM tags t
JOIN url_tags ut ON ut.tag_id = t.id
WHERE ut.url_id = $1; 