
## SETUP

1. Create redis.env inside docker/containers-configs
Example env for redis
```env
ALLOW_EMPTY_PASSWORD=yes
REDIS_PASSWORD=redispass
REDIS_AOF_ENABLED=no
```
2. Create secrets.env.yml inside docker/containers-configs/app-configs
``` yaml
db_config:
  password: "YourPasswordForAccountsService" # password, that you provided in postgre.env in accounts_service role (for encrypted password actual password, not hash)
  
crypto:
  bcrypt_cost: 5 # min 4, max 31

JWT:  
  change_password_token:
    TTL: 2h
    secret: "AnyString" # Any string for jwt tokens sign  
  verify_account_token:
    TTL: 3h
    secret: "AnyString" # Any string for jwt tokens sign  

redis_registration_options:
  password: "redispass" # Here is your password for redis with registration cache db 

session_cache_options:
  password: "redispass"  # Here is your password for redis with session cache db

account_sessions_cache_options:
  password: "redispass" # Here is your password for redis with account sessions cache db 
``` 