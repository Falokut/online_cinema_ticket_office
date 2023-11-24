## Configuration

1. Create .env in root dir and provide EMAIL_PASSWORD  
Example .env
```env
EMAIL_PASSWORD="passwordOrAPIKEY_ForEmail"
```
2. Expose this vars inside config.yml in folder  docker/containers-configs/app-configs
``` yaml
mail_sender:
  email_port: 465              # smtp port       
  email_host: "smtp.yandex.ru" # smtp host
  email_address: "Email"       # email of the sender of the emails
  email_login: "YourLogin"     # login for email
```