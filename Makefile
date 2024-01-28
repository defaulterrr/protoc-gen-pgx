build:
	go build -o bin/protoc-gen-pgx ./cmd/protoc-gen-pgx

test: 
	protoc --plugin=protoc-gen-pgx=bin/protoc-gen-pgx --pgx_opt=paths=source_relative --pgx_out=example -I example example/proto/type.proto