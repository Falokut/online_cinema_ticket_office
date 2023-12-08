CREATE ROLE movies_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title_ru TEXT NOT NULL,
    title_en TEXT,
    budget TEXT,
    plot TEXT NOT NULL,
    directors  INT ARRAY,
    genres  INT ARRAY,
    countries  INT ARRAY, 
    duration INT NOT NULL,
    poster_picture_id TEXT,
    cast_id INT NOT NULL,
    release_year SMALLINT NOT NULL
);
GRANT SELECT ON movies TO movies_service;