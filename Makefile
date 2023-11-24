project_name = gateway
api_version = v1

swagger-docs-dir = swagger/docs
		
create-swagger-dir:
	IF NOT EXIST "$(swagger-docs-dir)" ( MD "$(swagger-docs-dir)" )

swagger-doc-gen:
	echo("Hello")

.swagger:	create-swagger-dir	swagger-doc-gen	

.docker-compose:
	docker-compose --parallel -1 -f $(project_name).yml -p $(project_name) up --build $(project_name) -d --remove-orphans

.docker-compose-up:
	docker-compose --parallel -f $(project_name).yml up $(project_name) -d

.clear-images:
	docker image prune -f

.run:	.docker-compose		.clear-images