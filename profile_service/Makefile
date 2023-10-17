project_name = profile_service

gen:
	protoc --go_out=app/pkg --go-grpc_out=app/pkg app/protos/profilesService.proto
	
clean:
	rd pkg/$(project_name)/proto/*.pb.go

gen-cert:
	call app/cert/gen.cmd

server:
	go run app/cmd/server/app.go

docker-build:
	docker build -t $(project_name) .

docker-compose-build:
	docker-compose up --build $(project_name)
	
docker-compose-run:
	docker-compose up $(project_name)