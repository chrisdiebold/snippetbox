-- Create a `snippets` table.
CREATE TABLE snippets (
    id        SERIAL PRIMARY KEY,
    title     VARCHAR(100) NOT NULL,
    content   TEXT NOT NULL,
    created   TIMESTAMPTZ NOT NULL,
    expires   TIMESTAMPTZ NOT NULL
);

-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON snippets(created);