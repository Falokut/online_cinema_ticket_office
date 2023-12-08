CREATE ROLE geo_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name_ru text NOT NULL,
    name_en text
);
CREATE TABLE cities (
    city_id SERIAL PRIMARY KEY,
    country_id INT NOT NULL,
    name_ru text NOT NULL,
    name_en text
)

ALTER TABLE cities
    ADD CONSTRAINT country_id_fkey  FOREIGN KEY  (country_id) REFERENCES countries(id)

GRANT SELECT ON countries TO geo_service;
GRANT SELECT ON cities TO geo_service;