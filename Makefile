build: test
	go build -o ./bin/freqsim -v

test:
	go test ./

proto:
	protoc -I=protobuf/ --go_out=protobuf/ freqsim.proto

.PHONEY: clean

clean:
	go clean
	rm -rf bin
	mkdir bin
