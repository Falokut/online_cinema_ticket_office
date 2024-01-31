CREATE TABLE professions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
CREATE TABLE casts (
    movie_id INT NOT NULL,
    person_id INT NOT NULL,
    profession_id INT REFERENCES professions(id) ON DELETE SET NULL ON UPDATE CASCADE,
    PRIMARY KEY(movie_id,person_id,profession_id)
);

CREATE TABLE casts_labels (
    movie_id INT PRIMARY KEY,
    label TEXT NOT NULL 
);


CREATE OR REPLACE FUNCTION remove_cast()
RETURNS TRIGGER
AS $$
BEGIN
    IF NOT EXISTS (SELECT * FROM casts WHERE movie_id = OLD.movie_id) THEN
        DELETE FROM casts_labels WHERE movie_id = OLD.movie_id;
    END IF;
    RETURN NEW;
END; $$
LANGUAGE PLPGSQL;


CREATE TRIGGER remove_cast_trigger
    AFTER DELETE ON casts
    FOR EACH ROW
    EXECUTE PROCEDURE remove_cast();


GRANT SELECT, UPDATE, INSERT, DELETE ON casts TO admin_casts_service;
GRANT SELECT, UPDATE, INSERT, DELETE ON casts_labels TO admin_casts_service;
GRANT USAGE, SELECT ON SEQUENCE  professions_id_seq TO admin_casts_service;
GRANT SELECT, UPDATE, INSERT, DELETE ON professions TO admin_casts_service;

GRANT SELECT ON casts TO casts_service;
GRANT SELECT ON professions TO casts_service;


