CREATE TABLE IF NOT EXISTS tickers
(
    id         SERIAL PRIMARY KEY,
    symbol     VARCHAR(10) NOT NULL,
    notes      TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
