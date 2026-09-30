-- Migration: 000003_create_few_shot_examples_table.up.sql
CREATE TABLE IF NOT EXISTS few_shot_examples (
    uuid VARCHAR(64) PRIMARY KEY,
    archetype VARCHAR(32) NOT NULL,
    input_context TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
