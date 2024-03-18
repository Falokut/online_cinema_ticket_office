CREATE TABLE profiles
(
    account_id uuid NOT NULL PRIMARY KEY,
    email text NOT NULL UNIQUE, 
    username text NOT NULL,
    profile_picture_id text,
    registration_date date NOT NULL
);

GRANT SELECT,INSERT,DELETE,UPDATE ON profiles TO profiles_service;