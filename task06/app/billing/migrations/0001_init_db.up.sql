CREATE TABLE balance (
    user_id INTEGER PRIMARY KEY,
    amount INTEGER CHECK (amount > 0)
);

CREATE TABLE transaction (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES balance(user_id),
    order_id INTEGER UNIQUE,
    time TIMESTAMP NOT NULL,
    type VARCHAR(16) CHECK (type in ('expense', 'top_up')),
    amount INTEGER,
    status VARCHAR(16) CHECK (status in ('accepted', 'declined'))
);
