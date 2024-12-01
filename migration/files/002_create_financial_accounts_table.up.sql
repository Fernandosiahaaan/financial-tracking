CREATE TABLE IF NOT EXISTS finance_accounts (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

INSERT INTO finance_accounts (name, description, created_at, updated_at, deleted_at) 
    VALUES ('test', 'example insert', NOW(), NOW(), NULL)
;

