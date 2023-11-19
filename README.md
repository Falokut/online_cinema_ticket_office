# Online Cinema Ticket Office

This repository has the code for an online cinema ticket office. You can buy movie tickets and get showtime information.  

----

## Installation and Running

To install and run the service, follow these steps:

1. Clone the repository to your local machine using the following command:
   ```shell
   git clone https://github.com/Falokut/online_cinema_ticket_office.git
   ```

2. Set up the required services. Each service has its own setup instructions:
   * [Setup accounts_db](/accounts_db/README.md#SETUP)
   * [Setup accounts_service](/accounts_service/README.md#SETUP)
   * [Setup email_service](/email_service/README.md#SETUP)
   * [Setup profiles_service](/profiles_service/README.md#SETUP)

3. Start the gateway by running the following command:
   ```shell
   docker-compose -f gateway.yml up --build gateway
   ```
   If you have the `make` utility installed, you can also use the following command:
   ```shell
   make .docker-compose
   ```

4. Once the gateway is successfully started, you can access the RestAPI endpoint at `http://localhost:80` and the gRPC endpoint at `http://localhost:81`.
     
## Metrics and Monitoring

We use Grafana, Prometheus, and Jaeger to collect and visualize application metrics. You can track the performance of the application using these tools.

### Endpoints  
* Grafana endpoint  http://localhost:3000  
* Prometheus endpoint  http://localhost:9090
* Jaeger UI endpoint http://localhost:16686


## Accounts and authentication
The cinema ticket api features a login system where users can securely log in via sessions. This system ensures that only approved users can perform actions with their accounts.

To create an account, users can register by providing their email and password. Once registered, users can log in to their accounts using their credentials. The system will generate a session token for the user, which they will use for authentication in future requests.

Users remain logged in until they manually log out or their session expires. This eliminates the need for users to repeatedly authenticate themselves for each request, providing a seamless experience.

Users can safely access the online cinema ticket office using their account information. Additionally, it's worth noting that passwords are encrypted and not stored in plain text. Instead, they are encrypted using modern encryption algorithm bcrypt. This provides an added layer of security, as even in the event of a data breach, it would be extremely difficult for malicious actors to recover and exploit these passwords.

When registering a new account, the entered passwords are securely encrypted before being stored in the database. This way, user passwords are protected from unauthorized access.


## Development

We implement this project in the Golang programming language. You can help with the project by adding new features, fixing bugs, or improving the code.

## Author

- [Falokut](https://github.com/Falokut) - Primary author of the project

## License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---
