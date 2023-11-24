# Content

+ [Accounts service](#accounts-service)
    + [Features](#features)
    + [Accounts and authentication](#accounts-and-authentication)
    + [Registration](#registration)
+ [Configuration](#configuration)
+ [Docs](#docs)
---------

# Accounts service

The Account Service is a robust and secure service that provides essential functionalities for user accounts management. It offers a seamless user experience with features such as registration, password reset, account confirmation, login, and authentication.

# Features

1. Registration: Users can create new accounts by providing their basic information, including email and password. The registration process ensures that only valid and unique email addresses are accepted.

2. Password Reset: In case users forget their passwords, the service allows them to initiate a password reset procedure. A secure link is sent to the user's registered email address, enabling them to set a new password and regain access to their account.

3. Account Confirmation: To enhance security and prevent abuse, newly registered users must confirm their email addresses. A confirmation link is sent to the provided email, and upon verification, the account is activated within the system.

4. Login: Once registered and confirmed, users can securely log in to their accounts using their email and password. The service utilizes robust authentication protocols to protect account information and ensure secure access.

5. Authentication: To enhance security and prevent unauthorized access, the service employs authentication methods such as session-based identification and client identification based on their IP addresses. If the IP address provided in the request does not match the one stored in the session cache, access will be denied. These security measures ensure the safeguarding of user accounts and help in protecting against unauthorized access.

The Account Service provides a reliable, efficient, and user-friendly solution for managing user accounts in web applications. With its comprehensive set of features, it ensures the security and integrity of user data, delivering a seamless login and account management experience.

## Accounts and authentication
The accounts service features a login system where users can securely log in via sessions. This system ensures that only approved users can perform actions with their accounts.

To create an account, users can register by providing their email and password. Once registered and confirmed emails, users can log in to their accounts using their credentials. The system will generate a session token for the user, which they will use for authentication in future requests.

Users remain logged in until they manually log out or their session expires. This eliminates the need for users to repeatedly authenticate themselves for each request, providing a seamless experience.

Users can safely access the services using their account information. Additionally, it's worth noting that passwords are encrypted and not stored in plain text. Instead, they are encrypted using encryption algorithm bcrypt. This provides an added layer of security, as even in the event of a data breach, it would be extremely difficult for malicious actors to recover and exploit these passwords.

When registering a new account, the entered passwords are securely encrypted before being stored in the database. This way, user passwords are protected from unauthorized access.

## Registration
During the registration process, an email confirmation link is sent to the user's provided email address (need another request). The user must click on this link to verify their account and activate it. Once the email is confirmed, the account information is securely transferred from the Redis cache to the main database.

Implementing this email verification step helps ensure that only legitimate users with valid email addresses can create accounts on the cinema ticket. It helps prevent potential abuse or unauthorized access by requiring users to verify their identities before gaining full access to the system.

---
# Configuration
1. Create .env in root dir  
Example env for redis:
```env
REDIS_PASSWORD=redispass
REDIS_AOF_ENABLED=no
```
2. [Configure accounts_db](../databases/accounts_db/README.md#Configuration)
3. Create secrets.env.yml inside docker/containers-configs/app-configs  
Example config for service:
``` yaml
db_config:
  password: "YourPasswordForAccountsService" # password, that you provided in .env in accounts_db service
                                             # in accounts_service role (for encrypted password actual password, not hash)
  
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

# Docs
[Swagger docs](swagger/docs/accounts_service_v1.swagger.json)