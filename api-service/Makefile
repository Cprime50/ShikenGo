server:
	go run main.go


proto:
	rm -rf pb/*.go
	protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --proto_path=pb \
    pb/*.proto

path:
	PATH="${PATH}:${HOME}/go/bin"

key:
	chmod +x cert/gen.sh
	cert/gen.sh

test:
	cd src && go test -v
	cd utils && go test -v