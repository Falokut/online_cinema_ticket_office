## SETUP

1.  Create postgres.env inside containers-configs
Example env for postgre
```env
POSTGRES_USER=postgres
PGUSER=postgres
POSTGRES_PASSWORD=YourPassword
POSTGRES_DB="accounts_db"
PGDATA=/var/lib/postgresql/data
```
	
2. Change volume location for data inside accounts_db.yml, right part(after '':'') is PGDATA variable value
``` yaml
/.container_data/database/postgres/data:/var/lib/postgresql/data
```
3. Change passwords for roles inside init-up.sql in db folder
Your roles would look like this:
```sql
CREATE ROLE accounts_service WITH
    LOGIN
    CONNECTION LIMIT 100
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXMGszsnnHexVwOU=';  -- Here your password for accounts service

CREATE ROLE profiles_service WITH
    LOGIN
    CONNECTION LIMIT 100
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXMGszsnnHexVwOU='; -- Here your password for profiles service
```
or like this (passwords without encryption)
```sql
CREATE ROLE accounts_service WITH
    LOGIN
    CONNECTION LIMIT 100
    PASSWORD 'YourPasswordForAccountsService'; -- Here your password for accounts service


CREATE ROLE profiles_service WITH
    LOGIN
    CONNECTION LIMIT 100
    PASSWORD 'YourPasswordForProfilesService'; -- Here your password for profiles service
```