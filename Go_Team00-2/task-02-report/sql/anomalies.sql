CREATE TABLE anomalies (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(36),
    frequency DOUBLE PRECISION,
    timestamp TIMESTAMP
);
