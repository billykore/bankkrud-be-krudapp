CREATE TABLE IF NOT EXISTS pockets
(
    id             SERIAL PRIMARY KEY,
    uuid           UUID         NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id        UUID         NOT NULL,
    name           VARCHAR(255) NOT NULL,
    account_number VARCHAR(255) NOT NULL,
    saving_type    INTEGER      NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP WITH TIME ZONE     DEFAULT NULL
)