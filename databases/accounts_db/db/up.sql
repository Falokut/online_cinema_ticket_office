CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
    VERSION "1.1";

CREATE TABLE accounts
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    email text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    registration_date date NOT NULL DEFAULT now(),
    CONSTRAINT account_id_pkey PRIMARY KEY (id)
);
GRANT SELECT,DELETE,UPDATE,INSERT ON accounts TO accounts_service;