CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    email    VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR     NOT NULL
);
CREATE INDEX IF NOT EXISTS email_idx ON users (email);
