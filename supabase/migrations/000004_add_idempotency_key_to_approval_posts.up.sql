-- Migration: 000004_add_idempotency_key_to_approval_posts.up.sql
ALTER TABLE approval_posts ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(128) UNIQUE;
