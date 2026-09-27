-- +goose Up
INSERT INTO topics (name) VALUES ('BBŽ'), ('Hrvatska'), ('Svijet'), ('Tech');

-- +goose Down
DELETE FROM topics WHERE name IN ('BBŽ', 'Hrvatska', 'Svijet', 'Tech');
