-- Drop indexes
DROP INDEX IF EXISTS idx_urls_created_at;
DROP INDEX IF EXISTS idx_urls_user_id;
DROP INDEX IF EXISTS idx_urls_short_code;
 
-- Drop table
DROP TABLE IF EXISTS urls; 