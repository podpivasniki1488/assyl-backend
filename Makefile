

protoc:
	protoc -I proto proto/*.proto --go_out=./protopb --go_opt=paths=source_relative --go-grpc_out=./protopb --go-grpc_opt=paths=source_relative

swagger:
	swag init -g ./cmd/main.go
