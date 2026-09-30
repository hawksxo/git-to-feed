-- Migration: 000001_create_approval_posts_table.up.sql
CREATE TABLE IF NOT EXISTS approval_posts (
    uuid VARCHAR(64) PRIMARY KEY,
    content TEXT NOT NULL,
    archetype VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL,
    edited_content TEXT DEFAULT '',
    reject_reason TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);