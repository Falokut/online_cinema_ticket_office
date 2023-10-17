CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
    VERSION "1.1";
CREATE ROLE accounts_service WITH
    LOGIN
    CONNECTION LIMIT 100
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE ROLE profiles_service WITH
    LOGIN
    CONNECTION LIMIT 100
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE accounts
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    email text NOT NULL,
    password_hash text NOT NULL,
    registration_date date NOT NULL DEFAULT now(),
    CONSTRAINT account_id_pkey PRIMARY KEY (id)
);

CREATE TABLE profiles
(
    account_id uuid NOT NULL PRIMARY KEY,
    email text NOT NULL,
    username text NOT NULL,
    profile_picture_id text,
    registration_date date NOT NULL,
    CONSTRAINT account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) MATCH SIMPLE
      ON UPDATE NO ACTION
      ON DELETE CASCADE
);

GRANT ALL PRIVILEGES ON accounts TO accounts_service;
GRANT ALL PRIVILEGES ON profiles TO accounts_service;
GRANT ALL PRIVILEGES ON profiles TO profiles_service;