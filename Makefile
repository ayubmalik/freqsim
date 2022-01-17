build: test
	go build -o ./bin/server ./cmd/server
	go build -o ./bin/client ./cmd/client

test:
	go test ./...

proto:
	protoc -I=pb --go_out=pb --go-grpc_out=pb \
	  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative freqsim.proto

.PHONEY: clean

clean:
	go clean
	rm -rf bin
	mkdir bin
