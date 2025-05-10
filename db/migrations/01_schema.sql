-- Initialize the users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    CONSTRAINT users_uc_email UNIQUE (email)
);

-- Initialize the snippets table
CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Add an index to improve query performance for the Latest() method
CREATE INDEX idx_snippets_expires ON snippets(expires);

-- Note: If you want to add a foreign key relationship between users and snippets later:
-- ALTER TABLE snippets ADD COLUMN user_id INTEGER;
-- ALTER TABLE snippets ADD CONSTRAINT fk_snippets_users FOREIGN KEY (user_id) REFERENCES users(id);