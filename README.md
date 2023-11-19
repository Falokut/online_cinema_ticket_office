# Online Cinema Ticket Office

This repository contains the code for an online cinema ticket office, where you can buy movie tickets and get information about showtimes.

## Installation and Running

1. Clone the repository to your local machine:
   ```shell
   git clone https://github.com/Falokut/online_cinema_ticket_office.git
   ```
2. Setup all services:
* [Setup accounts_db](/accounts_db/README.md#SETUP)
* [Setup accounts_service](/accounts_service/README.md#SETUP)
* [Setup email_service](/email_service/README.md#SETUP)
* [Setup profiles_service](/profiles_service/README.md#SETUP)
1. Start gateway:
   ```shell
    docker-compose -f gateway.yml  up --build gateway
   ```
   or 
   ```shell
    make .docker-compose
   ```

## Metrics and Monitoring

Grafana and Prometheus are used to collect and visualize application metrics. You can track the performance of the application using these tools.

## Development

This project is implemented in the Golang programming language. You can contribute to the project by creating new features, fixing bugs, or improving existing code.

## Author

- [Falokut](https://github.com/Falokut) - Primary author of the project

## License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---
