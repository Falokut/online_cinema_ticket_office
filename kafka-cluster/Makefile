service_name = kafka


docker-build:
	docker build -t $(service_name) .

docker-compose-build:
	docker-compose up --build $(service_name)
	
docker-compose-up:
	docker-compose up $(service_name)