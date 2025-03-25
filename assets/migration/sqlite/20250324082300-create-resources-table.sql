-- +migrate Up
CREATE TABLE resources (
                           id TEXT PRIMARY KEY,
                           slug TEXT,
                           name TEXT,
                           description TEXT,
                           label TEXT,
                           type TEXT,
                           uri TEXT,
                           created_by TEXT,
                           updated_by TEXT,
                           created_at TIMESTAMP,
                           updated_at TIMESTAMP
);

-- +migrate Down
DROP TABLE resources;
