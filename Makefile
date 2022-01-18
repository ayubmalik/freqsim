build: test
	go build -o ./bin/server ./cmd/server
	go build -o ./bin/client ./cmd/client

test:
	go test ./...

proto:
	protoc -I=pb --go_out=pb --go-grpc_out=pb \
	  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative freqsim.proto

gen-cert:
	scripts/gen-cert.sh
	mv scripts/server.crt certs
	mv scripts/server.key certs
	mv scripts/ca.crt certs

.PHONEY: clean

clean:
	go clean
	rm -rf bin
	mkdir bin
