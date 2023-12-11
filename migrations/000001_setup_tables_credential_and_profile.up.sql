-- Enumerations
CREATE TYPE user_credential_type AS ENUM ('email', 'phone', 'username');
CREATE TYPE user_credential_status AS ENUM ('registered', 'active');

-- Table for user credentials
CREATE TABLE user_credentials (
    uid                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    credential_type    user_credential_type NOT NULL DEFAULT 'email',
    credential_access  VARCHAR(255) UNIQUE,
    credential_key     TEXT, -- Stored as bcrypt hash
    credential_salt    TEXT, -- Stored securely, see suggestion below
    status             user_credential_status NOT NULL DEFAULT 'registered',
    create_time        TIMESTAMPTZ DEFAULT current_timestamp,
    update_time        TIMESTAMPTZ,
    deleted_time       TIMESTAMPTZ,
    created_by         UUID,
    updated_by         UUID,
    deleted_by         UUID
);

-- Table for user profiles
CREATE TABLE user_profiles (
    uid               UUID PRIMARY KEY,
    username          VARCHAR(255) UNIQUE,
    first_name        VARCHAR(255),
    last_name         VARCHAR(255),
    thumbnail         VARCHAR(255),
    is_email_verified BOOLEAN DEFAULT false,
    create_time       TIMESTAMPTZ DEFAULT current_timestamp,
    update_time       TIMESTAMPTZ,
    deleted_time      TIMESTAMPTZ,
    created_by        UUID,
    updated_by        UUID,
    deleted_by        UUID,
    FOREIGN KEY (uid) REFERENCES user_credentials(uid) ON DELETE CASCADE
);
