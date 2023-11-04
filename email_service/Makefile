project_name = email_service

gen-cert:
	call app/cert/gen.cmd

server:
	go run app/cmd/server/app.go

docker-build:
	docker build -t $(project_name) .

.docker-compose:
	docker-compose up --build $(project_name)
	
docker-compose-run:
	docker-compose up $(project_name)