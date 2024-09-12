-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE user_sessions (
    token_hash VARCHAR(512) NOT NULL,
    user_id UUID REFERENCES users(id),
    access_token_id UUID NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    UNIQUE(token_hash, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE user_sessions;
-- +goose StatementEnd
