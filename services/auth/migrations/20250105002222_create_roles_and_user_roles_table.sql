-- migrate:up

CREATE TABLE IF NOT EXISTS roles(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE users
ADD COLUMN role_id INT NOT NULL;

ALTER TABLE users
ADD CONSTRAINT fk_users_roles FOREIGN KEY (role_id) REFERENCES roles(id)
ON DELETE CASCADE;

-- migrate:down

DROP TABLE IF EXISTS roles;

ALTER TABLE users
DROP fk_users_roles;

ALTER TABLE users
DROP COLUMN role_id;