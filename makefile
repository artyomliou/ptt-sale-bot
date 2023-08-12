all: pttsalebot_client pttsalebot_server

pttsalebot_client: client/main.go
	cd client && CGO_ENABLED=0 go build . && mv client ../pttsalebot_client

pttsalebot_server: server/main.go
	cd server && CGO_ENABLED=0 go build . && mv server ../pttsalebot_server

grpc_api/%.pb.go: $(wildcard grpc_api/*.proto)
	protoc --proto_path=grpc_api --go_out=grpc_api --go_opt=paths=source_relative --go-grpc_out=grpc_api --go-grpc_opt=paths=source_relative $?

clean:
	rm $(wildcard grpc_api/*.pb.go)
	rm pttsalebot