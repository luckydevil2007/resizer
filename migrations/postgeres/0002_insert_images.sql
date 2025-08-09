-- +goose Up
INSERT INTO notes (note_title, owner_id, note_path, lat, lon ) VALUES ('1st note', 1, 'storage/1.ogg',59.9371, 30.4263 );
-- +goose Down

