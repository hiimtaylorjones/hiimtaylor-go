-- +goose Up
ALTER TABLE posts ADD COLUMN banner_image_url VARCHAR(500);

-- +goose Down
ALTER TABLE posts DROP COLUMN banner_image_url;
