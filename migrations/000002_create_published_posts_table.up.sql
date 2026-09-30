-- Migration: 000002_create_published_posts_table.up.sql
CREATE TABLE IF NOT EXISTS published_posts (
    id VARCHAR(64) PRIMARY KEY,
    approval_post_uuid VARCHAR(64) NOT NULL REFERENCES approval_posts(uuid),
    linkedin_share_urn VARCHAR(128) DEFAULT '',
    status VARCHAR(32) NOT NULL,
    error_message TEXT DEFAULT '',
    published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);