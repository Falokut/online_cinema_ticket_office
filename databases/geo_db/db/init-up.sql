CREATE ROLE geo_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE ROLE admin_geo_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name_ru text NOT NULL,
    name_en text
);

CREATE TABLE cities (
    city_id SERIAL PRIMARY KEY,
    country_id INT REFERENCES countries(id) ON DELETE SET NULL,
    name_ru text NOT NULL,
    name_en text
);

GRANT SELECT, UPDATE, DELETE, INSERT ON countries TO admin_geo_service;
GRANT SELECT, UPDATE, DELETE, INSERT ON cities TO admin_geo_service;
GRANT USAGE, SELECT ON SEQUENCE cities_city_id_seq TO admin_geo_service;
GRANT USAGE, SELECT ON SEQUENCE countries_id_seq TO admin_geo_service;
GRANT SELECT ON countries TO geo_service;
GRANT SELECT ON cities TO geo_service;