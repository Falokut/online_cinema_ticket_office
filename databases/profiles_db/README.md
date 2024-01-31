# Configuration
1.  Create .env in root dir
Example env for postgres
```env
POSTGRES_USER=postgres
PGUSER=postgres
POSTGRES_PASSWORD=YourPassword
SERVICE_PASSWORD=YourPassword
```
2. setup pgbouncer:
* create userlist.txt in docker/pgbouncer and provide passwords: 
```
"profiles_service" "yourpassword"
"postgres" "yourpassword"
```