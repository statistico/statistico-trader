build:
	protoc --proto_path=. --go_out=plugins=grpc:. internal/app/grpc/proto/*.proto