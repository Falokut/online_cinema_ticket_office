project_name = online_cinema_ticket_office
service_to_up = gateway
api_version = v1

.docker-compose:
	docker-compose --parallel -1 -f $(service_to_up).yml -p $(project_name) up --build $(service_to_up) -d --remove-orphans

.docker-compose-up:
	docker-compose --parallel -1 -f $(service_to_up).yml up $(service_to_up) -d

.clear-images:
	docker image prune -f

.run:	.docker-compose		.clear-images