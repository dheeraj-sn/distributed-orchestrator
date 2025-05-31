CREATE TABLE job_results (
    job_id UUID PRIMARY KEY REFERENCES jobs(id) ON DELETE CASCADE,
    output TEXT,
    logs TEXT,
    created_at TIMESTAMP DEFAULT now()
);