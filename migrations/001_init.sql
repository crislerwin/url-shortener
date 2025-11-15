-- Create urls table
CREATE TABLE IF NOT EXISTS urls (
    short_id VARCHAR(6) PRIMARY KEY,
    encrypted_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    click_count BIGINT DEFAULT 0,
    last_accessed TIMESTAMP WITH TIME ZONE
);

-- Create index for expiration queries
CREATE INDEX IF NOT EXISTS idx_expires_at ON urls(expires_at) WHERE expires_at IS NOT NULL;

-- Create index for created_at for analytics
CREATE INDEX IF NOT EXISTS idx_created_at ON urls(created_at);

-- Add comment to table
COMMENT ON TABLE urls IS 'Stores shortened URLs with encryption';
COMMENT ON COLUMN urls.short_id IS 'Unique 6-character alphanumeric identifier';
COMMENT ON COLUMN urls.encrypted_url IS 'AES-256 encrypted original URL';
COMMENT ON COLUMN urls.created_at IS 'Timestamp when the short URL was created';
COMMENT ON COLUMN urls.expires_at IS 'Optional expiration timestamp for the URL';
COMMENT ON COLUMN urls.click_count IS 'Number of times the short URL has been accessed';
COMMENT ON COLUMN urls.last_accessed IS 'Timestamp of the last access';
