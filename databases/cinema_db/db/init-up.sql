CREATE ROLE cinema_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';



CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE cinemas (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    city_id INT REFERENCES cities(id) ON UPDATE CASCADE ON DELETE SET NULL,
    --address without city
    address TEXT NOT NULL,
    coordinates geography(POINT,4326) NOT NULL
);

CREATE TABLE halls_types (
    type_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE halls (
    id SERIAL PRIMARY KEY,
    cinema_id INT REFERENCES cinemas(id) ON UPDATE CASCADE ON DELETE SET NULL,
    hall_type_id INT REFERENCES halls_types(type_id) ON UPDATE CASCADE ON DELETE SET NULL,
    name TEXT NOT NULL,
    hall_size INT NOT NULL DEFAULT 0
);

CREATE TABLE halls_configurations (
    hall_id INT REFERENCES halls(id) ON UPDATE CASCADE,
    row INT CHECK(row > 0),
    seat INT CHECK(seat > 0),
    grid_pos_x FLOAT NOT NULL,
    grid_pos_y FLOAT NOT NULL,
    PRIMARY KEY(hall_id, row, seat)
);

CREATE OR REPLACE FUNCTION update_hall_size()
RETURNS TRIGGER
AS $$
BEGIN
    UPDATE halls SET hall_size=(SELECT COUNT(seat) FROM halls_configurations WHERE hall_id=id);
    RETURN NEW;
END; $$
LANGUAGE PLPGSQL;

CREATE TRIGGER hall_place_insert_trigger
            AFTER INSERT ON halls_configurations
            EXECUTE FUNCTION update_hall_size();

CREATE TRIGGER hall_place_delete_trigger
            AFTER DELETE ON halls_configurations
            FOR EACH ROW
            EXECUTE FUNCTION update_hall_size();


CREATE TABLE screenings_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE screenings (
    id BIGSERIAL PRIMARY KEY,
    screening_type_id INT REFERENCES screenings_types(id) ON UPDATE CASCADE ON DELETE SET NULL,
    movie_id INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL CHECK(start_time > clock_timestamp()),
    hall_id INT REFERENCES halls(id) ON UPDATE CASCADE,
    ticket_price DECIMAL(8,2) CHECK(ticket_price>0.0)
);


GRANT SELECT ON cities TO cinema_service;
GRANT SELECT ON cinemas TO cinema_service;
GRANT SELECT ON halls_configurations TO cinema_service;
GRANT SELECT ON halls_types TO cinema_service;

GRANT SELECT ON halls TO cinema_service;
GRANT SELECT ON screenings TO cinema_service;
GRANT SELECT ON screenings_types TO cinema_service;





