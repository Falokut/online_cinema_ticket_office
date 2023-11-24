## Configuration
1.  Create .env in root dir
Example env for postgre
```env
POSTGRES_USER=postgres
PGUSER=postgres
POSTGRES_PASSWORD=YourPassword
```	
2. Change passwords for roles inside init-up.sql in db folder
Your role would look like this:
```sql
CREATE ROLE profiles_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXMGszsnnHexVwOU=';  -- Here your password for profiles service
```
or like this (passwords without encryption)
```sql
CREATE ROLE profiles_service WITH
    LOGIN
    PASSWORD 'YourPasswordForProfilesService'; -- Here your password for profiles service
```