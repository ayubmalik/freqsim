build: test
	go build -o ./bin/freqsim -v

test:
	go test ./

proto:
	protoc -I=protobuf --go_out=protobuf --go-grpc_out=protobuf \
	  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative freqsim.proto

.PHONEY: clean

clean:
	go clean
	rm -rf bin
	mkdir bin
