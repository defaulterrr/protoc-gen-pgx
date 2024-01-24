build:
	go build -o bin/protoc-gen-pgx ./cmd/protoc-gen-pgx

test: 
	protoc --plugin=protoc-gen-pgx=bin/protoc-gen-pgx --pgx_out=. -I example example/proto/type.proto