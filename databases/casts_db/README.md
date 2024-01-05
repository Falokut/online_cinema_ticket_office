# Configuration
1.  Create .env in root dir
Example env for postgre
```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=YourPassword
```	
2. Change passwords for roles inside init-up.sql in db folder
Your role would look like this:
```sql
CREATE ROLE casts_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXMGszsnnHexVwOU=';  -- Here your password for casts service
CREATE ROLE admin_casts_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXMGszsnnHexVwOU=';  -- Here your password for admin casts service
```
or like this (passwords without encryption)
```sql
CREATE ROLE casts_service WITH
    LOGIN
    PASSWORD 'YourPasswordForCastsService'; -- Here your password for casts service

CREATE ROLE admin_casts_service WITH
    LOGIN
    PASSWORD 'YourPasswordForCastsService'; -- Here your password for admin casts service
```

3. setup pgbouncer:
* create userlist.txt in docker/pgbouncer and provide passwords: 
```
"casts_service" "yourpassword"
"admin_casts_service" "yourpassword"
"postgres" "yourpassword"
```