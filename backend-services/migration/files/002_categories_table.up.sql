CREATE TYPE transaction_type AS ENUM ('INCOME', 'EXPENSE', 'TRANSFER');

-- Tabel categories
CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    type transaction_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories_hist (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    type transaction_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);