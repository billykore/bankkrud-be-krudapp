CREATE TABLE users
(
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id       UUID UNIQUE              DEFAULT gen_random_uuid(),
    username      VARCHAR(50) UNIQUE  NOT NULL,
    first_name    VARCHAR(100)        NOT NULL,
    last_name     VARCHAR(100)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    phone_number  VARCHAR(20) UNIQUE  NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    date_of_birth DATE                NOT NULL,
    cif_number    VARCHAR(20) UNIQUE  NOT NULL,
    last_login    TIMESTAMP,
    status        VARCHAR(20)              DEFAULT 'active',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    address       TEXT       NOT NULL,
    cif           VARCHAR(20) UNIQUE  NOT NULL
);
