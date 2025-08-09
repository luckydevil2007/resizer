-- +goose Up

CREATE TABLE
users(
    id SERIAL PRIMARY KEY,
    user_name  text,
    password_hash text
);


CREATE TABLE notes (
    id SERIAL NOT NULL PRIMARY KEY,
    note_title TEXT,
    owner_id INT,
    note_path TEXT,
    lat FLOAT, 
    lon FLOAT
);


INSERT INTO users (user_name, password_hash) VALUES ('user','123');

-- +goose Down
DROP TABLE images;
DROP TABLE users;