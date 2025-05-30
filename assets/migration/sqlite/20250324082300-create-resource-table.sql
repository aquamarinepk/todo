-- +migrate Up
CREATE TABLE resource (
                           id TEXT PRIMARY KEY,
                           short_id TEXT,
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
DROP TABLE resource;
