CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
    bank_id bigint NOT NULL REFERENCES banks ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (bank_id, permission_id)
);

-- Add the two permissions to the table.
-- INSERT INTO permissions (code)
-- VALUES 
--     ('movies:read'),
--     ('movies:write');