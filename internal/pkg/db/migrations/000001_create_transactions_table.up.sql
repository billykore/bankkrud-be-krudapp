CREATE TABLE IF NOT EXISTS transactions
(
    id                    SERIAL PRIMARY KEY,
    uuid                  UUID           NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id               UUID           NOT NULL UNIQUE,
    pocket_id             UUID           NOT NULL UNIQUE,
    destination_account   VARCHAR(255)   NOT NULL,
    transaction_type      VARCHAR(255)   NOT NULL,
    transaction_reference VARCHAR(255)   NOT NULL,
    status                VARCHAR(255)   NOT NULL,
    note                  VARCHAR(255),
    amount                NUMERIC(10, 2) NOT NULL,
    fee                   NUMERIC(10, 2) NOT NULL        DEFAULT 0,
    created_at            TIMESTAMP WITH TIME ZONE       DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP WITH TIME ZONE       DEFAULT CURRENT_TIMESTAMP,
    deleted_at            TIMESTAMP WITH TIME ZONE       DEFAULT NULL
);