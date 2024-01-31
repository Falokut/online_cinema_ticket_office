CREATE TABLE persons (
    id SERIAL PRIMARY KEY,
    fullname_ru TEXT NOT NULL,
    fullname_en TEXT,
    birthday DATE,
    sex TEXT,
    photo_id TEXT
);

GRANT SELECT ON persons TO movies_persons_service;
GRANT SELECT, UPDATE, DELETE, INSERT ON persons TO admin_movies_persons_service;
GRANT USAGE, SELECT ON SEQUENCE persons_id_seq TO admin_movies_persons_service;