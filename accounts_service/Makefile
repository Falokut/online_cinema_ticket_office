project_name = accounts_service
api_version = v1
protoc_out_dir = app/pkg/$(project_name)/$(api_version)/protos

proto_files_dir = app/proto/$(project_name)/$(api_version)

protoc-gen:
	protoc  \
	-I include/googleapis -I include/grpc-gateway \
	--go_opt=M$(proto_files_dir)/api.proto=$(protoc_out_dir) \
	--go_opt=M$(proto_files_dir)/messages.proto=$(protoc_out_dir) \
	--go_out=app/pkg --go-grpc_out=app/pkg \
   	api.proto messages.proto -I $(proto_files_dir)

gateway-gen:
	protoc -I include/googleapis -I include/grpc-gateway \
	--grpc-gateway_out=logtostderr=true,paths=source_relative:./$(protoc_out_dir) \
   	api.proto messages.proto -I $(proto_files_dir)


swagger-docs-dir = swagger/docs
		
create-swagger-dir:
	IF NOT EXIST "$(swagger-docs-dir)" ( MD "$(swagger-docs-dir)" )

swagger-doc-gen:
	protoc -I include/googleapis -I include/grpc-gateway \
	--openapiv2_out ./$(swagger-docs-dir) \
	api.proto -I $(proto_files_dir)

.swagger:	create-swagger-dir	swagger-doc-gen	
.protoc:	protoc-gen	gateway-gen	.swagger

gen-cert:
	call app/cert/gen.cmd

server:
	go run app/cmd/server/app.go


.docker-build:
	docker build -t $(project_name) .

.docker-compose:
	docker-compose -f $(project_name).yml -p $(project_name) up --build $(project_name)
	
.docker-compose-up:
	docker-compose up $(project_name)