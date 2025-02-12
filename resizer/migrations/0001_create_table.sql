-- +goose Up
CREATE TABLE
users(
    id int not NULL,
    user_name  text
)

CREATE TABLE images (
    id int NOT NULL, serial=int not null auto increment
    image_title text,
    owner_id int,
    image_path text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE images;
DROP TABLE users;