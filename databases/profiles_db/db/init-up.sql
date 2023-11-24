CREATE ROLE profiles_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';


CREATE TABLE profiles
(
    account_id uuid NOT NULL PRIMARY KEY,
    email text NOT NULL,
    username text NOT NULL,
    profile_picture_id text,
    registration_date date NOT NULL
);
GRANT ALL PRIVILEGES ON profiles TO profiles_service;