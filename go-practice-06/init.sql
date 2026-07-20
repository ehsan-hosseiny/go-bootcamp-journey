CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price INT NOT NULL,
    stock INT NOT NULL
);

INSERT INTO products (name, price, stock)
VALUES
    ('PS5 Slim', 1000, 5),
    ('Keyboard', 120, 10),
    ('Mouse', 50, 20)
ON CONFLICT (name) DO NOTHING;
