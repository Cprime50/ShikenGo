server:
	go run main.go


proto:
	rm -rf profilepb/*.go
	protoc --go_out=profilepb --go_opt=paths=source_relative \
    --go-grpc_out=profilepb --go-grpc_opt=paths=source_relative \
    --proto_path=profilepb \
    profilepb/*.proto

key:
	chmod +x cert/gen.sh
	cert/gen.sh

test:
	cd src && go test -v
	cd utils && go test -v

path:
	PATH="${PATH}:${HOME}/go/bin"