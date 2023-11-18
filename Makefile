project_name = gateway
api_version = v1

swagger-docs-dir = swagger/docs
		
create-swagger-dir:
	IF NOT EXIST "$(swagger-docs-dir)" ( MD "$(swagger-docs-dir)" )

swagger-doc-gen:
	echo("Hello")

.swagger:	create-swagger-dir	swagger-doc-gen	

gen-cert:
	call app/cert/gen.cmd

server:
	go run app/cmd/server/app.go

.docker-build:
	docker build -t $(project_name) .

.docker-compose:
	docker-compose -f $(project_name).yml -p $(project_name) up --build $(project_name)
	
.docker-compose-up:
	docker-compose -f $(project_name).yml up $(project_name)