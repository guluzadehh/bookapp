-- migrate:up

CREATE INDEX IF NOT EXISTS idx_users_email 
ON users (email);

-- migrate:down

DROP INDEX IF EXISTS idx_users_email;