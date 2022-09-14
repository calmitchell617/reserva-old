CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    bank_id bigint NOT NULL REFERENCES banks ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);