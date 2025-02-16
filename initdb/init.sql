CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
    coins INTEGER NOT NULL DEFAULT 1000
    );

CREATE TABLE IF NOT EXISTS inventory (
                                         id SERIAL PRIMARY KEY,
                                         username VARCHAR(255) NOT NULL,
    item VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    UNIQUE(username, item)
    );

CREATE TABLE IF NOT EXISTS coin_transactions (
                                                 id SERIAL PRIMARY KEY,
                                                 from_user VARCHAR(255),
    to_user VARCHAR(255),
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );