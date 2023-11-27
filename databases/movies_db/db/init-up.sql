CREATE ROLE movies_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    genre_name TEXT NOT NULL
);

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    plot TEXT NOT NULL,
    poster_picture_id TEXT,
    cast_id INT
);

CREATE TABLE movies_genres (
    movie_id INT,
    genre_id INT,
    FOREIGN KEY (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres (id) ON DELETE CASCADE
);

GRANT ALL PRIVILEGES ON movies TO movies_service;
GRANT ALL PRIVILEGES ON movies TO movies_service;
GRANT ALL PRIVILEGES ON movies TO movies_service;