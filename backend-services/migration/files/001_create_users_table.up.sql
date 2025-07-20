CREATE TYPE role_type AS ENUM ('USER', 'ADMIN', 'SUPERADMIN');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    phone_number VARCHAR(100) UNIQUE NOT NULL,
    role role_type NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_hist (
    id UUID PRIMARY KEY UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    phone_number VARCHAR(100) UNIQUE NOT NULL,
    role role_type NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (
    -- id,
    username,
    password,
    full_name,
    email,
    phone_number,
    role
) VALUES (
    -- '550e8400-e29b-41d4-a716-446655440000',  -- UUID valid
    'user_test',
    'hashed_password_here',
    'test',
    'test@example.com',
    '+628123456789',
    'ADMIN'  -- Asumsikan 'admin' adalah salah satu nilai valid dari ENUM role_type
);
