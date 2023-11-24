project_name = profiles_service
api_version = v1
protoc_out_dir = app/pkg/$(project_name)/$(api_version)/protos
protoc_out_dir_win = app\pkg\$(project_name)\$(api_version)\protos
proto_files_dir = app/proto/$(project_name)/$(api_version)

protoc-clear:
	IF EXIST "$(protoc_out_dir_win)" (rd $(protoc_out_dir_win) /q /s)

protoc-gen:
	protoc  \
	-I include/googleapis -I include/grpc-gateway \
	--go_opt=M$(proto_files_dir)/$(project_name)_$(api_version).proto=$(protoc_out_dir) \
	--go_opt=M$(proto_files_dir)/$(project_name)_$(api_version)_messages.proto=$(protoc_out_dir) \
	--go_out=app/pkg --go-grpc_out=app/pkg \
   	$(project_name)_$(api_version).proto $(project_name)_$(api_version)_messages.proto -I $(proto_files_dir)

gateway-gen:
	protoc -I include/googleapis -I include/grpc-gateway \
	--grpc-gateway_out=logtostderr=true,paths=source_relative:./$(protoc_out_dir) \
   	$(project_name)_$(api_version).proto $(project_name)_$(api_version)_messages.proto -I $(proto_files_dir)


swagger-docs-dir = swagger/docs
swagger-docs-dir-win = swagger\docs

swagger-clear:
	IF EXIST "$(swagger-docs-dir-win)\$(project_name)_$(api_version).swagger.json" (del $(swagger-docs-dir-win)\$(project_name)_$(api_version).swagger.json /q /s)

create-swagger-dir:
	IF NOT EXIST "$(swagger-docs-dir)" ( MD "$(swagger-docs-dir)" )

swagger-doc-gen:
	protoc -I include/googleapis -I include/grpc-gateway \
	--openapiv2_out ./$(swagger-docs-dir) \
	$(project_name)_$(api_version).proto -I $(proto_files_dir)

.swagger:	swagger-clear	create-swagger-dir	swagger-doc-gen	
.protoc:	protoc-clear	protoc-gen	gateway-gen	.swagger