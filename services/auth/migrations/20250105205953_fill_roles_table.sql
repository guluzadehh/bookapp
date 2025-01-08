-- migrate:up

INSERT INTO roles(name)
VALUES ('user'), ('admin');

-- migrate:down

DELETE FROM roles WHERE name = 'user' OR name = 'admin';
