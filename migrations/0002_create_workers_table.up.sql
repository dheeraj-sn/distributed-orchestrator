CREATE TABLE workers (
    id TEXT PRIMARY KEY,
    host TEXT NOT NULL,
    last_heartbeat TIMESTAMP NOT NULL DEFAULT now()
);