CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    hash VARCHAR(255) NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_urls_unique_hash ON urls(hash);