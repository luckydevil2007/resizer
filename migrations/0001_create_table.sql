-- +goose Up
CREATE TABLE
users(
    id SERIAL PRIMARY KEY,
    user_name  text,
    password_hash text
);

CREATE SEQUENCE images_id_seq;
CREATE TABLE images (
    id INT NOT NULL DEFAULT nextval('images_id_seq') PRIMARY KEY,
    image_title TEXT,
    owner_id INT,
    image_path TEXT
);

INSERT INTO users (user_name, password_hash) VALUES ('user','123');

-- +goose Down
DROP TABLE images;
DROP TABLE users;