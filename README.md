# Content

+ [Online Cinema Ticket Office](#online-cinema-ticket-office)
+ [Installation and Running](#installation-and-running)
   + [Instructions](#instructions-to-run-the-gateway)
   + [Docker troubleshooting](#checking-and-updating-docker-compose-version)
+ [Services](#services)
+ [Metrics and Monitoring](#metrics-and-monitoring)
+ [Endpoints](#endpoints)

# Online Cinema Ticket Office

This repository has the code for an online cinema ticket office. You can buy movie tickets and get showtime information.  

----

# Installation and Running

**Prerequisites**:
- Docker installed on your system.
- Docker Compose installed on your system.

## Instructions to run the Gateway:

1. Clone the repository to your local machine using the following command:
   ```shell
   git clone https://github.com/Falokut/online_cinema_ticket_office.git
   ```

2. Configure the required services. Each service has its own setup instructions:
   * Configure databases inside databases folder (and skip first item on service configuration, doesn't need configure databases twise)
   * [Configure accounts service](https://github.com/Falokut/accounts_service/blob/master/README.md#Configure)
   * [Configure email service](https://github.com/Falokut/email_service/blob/master/README.md#Configure)
   * [Configure profiles service](https://github.com/Falokut/profiles_service/blob/master/README.md#Configure)
   * [Configure movies service](https://github.com/Falokut/movies_service/blob/master/README.md#Configure)
   * [Configure movies persons service](https://github.com/Falokut/movies_persons_service/blob/master/README.md#Configure)
   * [Configure casts service](https://github.com/Falokut/casts_service/blob/master/README.md#Configure)
   * [Configure cinema service](https://github.com/Falokut/cinema_service/blob/master/README.md#Configure)
   * [Configure cinema orders service](https://github.com/Falokut/cinema_orders_service/blob/master/README.md#Configure)
   * [Configure admin movies persons service](https://github.com/Falokut/admin_movies_persons_service/blob/master/README.md#Configure)
   * [Configure admin movies service](https://github.com/Falokut/admin_movies_service/blob/master/README.md#Configure)
   * [Configure admin casts service](https://github.com/Falokut/admin_casts_service/blob/master/README.md#Configure)
   * [Configure admin cinema service](https://github.com/Falokut/admin_cinema_service/blob/master/README.md#Configure)

3. Start the gateway by running the following command:

   ```shell
   docker-compose -f gateway.yml up --build gateway
   ```
   If you have the `make` utility installed, you can also use the following command:
   ```shell
   make .docker-compose
   ```

4. Once the gateway is successfully started, you can access the RestAPI endpoint at `http://localhost:8000` and `https://localhost:44300`.
---

Please note that these instructions assume that Docker and Docker Compose are already installed on your system. If you haven't installed them yet, please refer to the Docker documentation for the appropriate installation steps for your operating system.

## Checking and updating Docker compose version
To check the version of Docker Compose, you can use the following command:

```shell
docker-compose --version
```
This will display the version number of Docker Compose installed on your system.

If you need to update Docker Compose to a newer version, follow these instuctions:
Sure! Here are the additional instructions for Windows and macOS:

**Instructions for Windows:**

1. Download the latest binary of Docker Compose by visiting the official GitHub release page: [https://github.com/docker/compose/releases](https://github.com/docker/compose/releases).

2. Scroll down to the "Assets" section and find the binary that matches your system architecture, typically the one ending with `.exe` (e.g., `docker-compose-Windows-x86_64.exe`).

3. Click on the binary to download it.

4. Move the downloaded binary to a directory in your system's `PATH` environment variable. This allows you to run Docker Compose from anywhere.

**Instructions for Linux and macOS:**

1. Open a terminal.

2. Download the latest binary of Docker Compose using the following command:
   ```shell
   sudo curl -L "https://github.com/docker/compose/releases/download/{VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   ```
   Replace `{VERSION}` with the actual version number you want to install.

3. Apply executable permissions to the Docker Compose binary:
   ```shell
   sudo chmod +x /usr/local/bin/docker-compose
   ```

4. Verify that Docker Compose has been updated successfully by running:
   ```shell
   docker-compose --version
   ```
   It should display the newly installed version.
---

# Services
   + [Accounts service](https://github.com/Falokut/accounts_service)
   + [Profiles service](https://github.com/Falokut/profiles_service)
   + [Images storage service](https://github.com/Falokut/images_storage_service)  
   + [Image processing service](https://github.com/Falokut/image_processing_service)
   + [Email service](https://github.com/Falokut/email_service)
   + [Cinema service](https://github.com/Falokut/cinema_service)
   + [Cinema orders service](https://github.com/Falokut/cinema_orders_service)
   + [Movies service](https://github.com/Falokut/movies_service)
   + [Casts service](https://github.com/Falokut/casts_service)
   + [Movies persons service](https://github.com/Falokut/movies_persons_service)
   + [Admin cinema service](https://github.com/Falokut/admin_cinema_service)
   + [Admin movies service](https://github.com/Falokut/admin_movies_service)
   + [Admin casts service](https://github.com/Falokut/admin_casts_service)
   + [Admin movies persons service](https://github.com/Falokut/admin_movies_persons_service)

# Endpoints
* Grafana endpoint  http://localhost:3000  
* Prometheus endpoint  http://localhost:9090
* Jaeger UI endpoint http://localhost:16686
* RestApi and gRPC users endpoint http://localhost:8000 or https://localhost:44300
* kafka-ui users http://localhost:18082
* kafka-ui admins/employee http://localhost:28082
* pgAdmin http://localhost:100
* RestApi and gRPC admins/employee endpoint http://localhost:10080 or https://localhost:10443

# Metrics and Monitoring

We use Grafana, Prometheus, and Jaeger to collect and visualize application metrics. You can track the performance of the application using these tools.

# Author

- [@Falokut](https://github.com/Falokut) - Primary author of the project

# License

This project is licensed under the terms of the [MIT License](https://opensource.org/licenses/MIT).

---
