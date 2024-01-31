# Configuration
1.  Create .env in root dir
Example env for postgres
```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=YourPassword
ADMIN_SERVICE_PASSWORD=YourPassword
SERVICE_PASSWORD=YourPassword
```

2. setup pgbouncer:
* create userlist.txt in docker/pgbouncer and provide passwords: 
```
"movies_persons_service" "yourpassword"
"admin_movies_persons_service" "yourpassword"
"postgres" "yourpassword"
```