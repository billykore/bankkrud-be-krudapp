CREATE TABLE IF NOT EXISTS schedules
(
    id                 SERIAL PRIMARY KEY,
    uuid               UUID           NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    username           VARCHAR(100)   NOT NULL,
    name               VARCHAR(255)   NOT NULL,
    transaction_type   VARCHAR(100)   NOT NULL,
    source_account     VARCHAR(100)   NOT NULL,
    destination_number VARCHAR(100)   NOT NULL,
    amount             NUMERIC(14, 0) NOT NULL,
    period             VARCHAR(100)   NOT NULL,
    status             VARCHAR(100)   NOT NULL,
    cron_expression    VARCHAR(20)    NOT NULL,
    created_at         TIMESTAMP WITH TIME ZONE       DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP WITH TIME ZONE       DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMP WITH TIME ZONE       DEFAULT NULL
);