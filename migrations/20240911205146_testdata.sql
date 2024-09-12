-- +goose Up
-- +goose StatementBegin
-- INSERT INTO users (id, email) 
-- VALUES (gen_random_uuid(), 'testuser@statebyte.dev');
INSERT INTO users (id, email) 
VALUES ('5f32a8a2-45f4-4c0c-9f7e-61e8b8d4edc9', 'testuser@statebyte.dev');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
