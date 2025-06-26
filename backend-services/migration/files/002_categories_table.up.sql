-- ENUM untuk tipe transaksiAdd commentMore actions
CREATE TYPE transaction_type AS ENUM ('income', 'expense');

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