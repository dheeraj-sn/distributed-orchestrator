CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    task TEXT NOT NULL,
    args TEXT[],
    status TEXT NOT NULL,
    result TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);