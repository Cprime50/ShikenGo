server:
	go run main.go


proto:
	rm -rf quizpb/*.go
	protoc --go_out=quizpb --go_opt=paths=source_relative \
    --go-grpc_out=quizpb --go-grpc_opt=paths=source_relative \
    --proto_path=quizpb \
    quizpb/*.proto

key:
	chmod +x cert/gen.sh
	cert/gen.sh

path:
	PATH="${PATH}:${HOME}/go/bin"