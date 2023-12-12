CREATE ROLE movies_persons_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE persons (
    id SERIAL PRIMARY KEY,
    fullname_ru TEXT NOT NULL,
    fullname_en TEXT,
    birthday DATE,
    sex TEXT,
    photo_id TEXT
);

GRANT SELECT ON persons TO movies_persons_service;