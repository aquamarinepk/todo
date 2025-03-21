-- Res: User
-- Table: users

-- GetAll
SELECT id, username, email, slug FROM users;

-- Get
SELECT id, username, email, slug FROM users WHERE id = ?;

-- Create
INSERT INTO users (id, username, email, slug) VALUES (?, ?, ?, ?);

-- Update
UPDATE users SET username = ?, email = ?, slug = ? WHERE id = ?;

-- Delete
DELETE FROM users WHERE id = ?;
