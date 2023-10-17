project_name = api_service

gen-cert:
	call app/cert/gen.cmd

server:
	go run app/cmd/server/app.go

build-docker:
	docker build -t $(project_name) .