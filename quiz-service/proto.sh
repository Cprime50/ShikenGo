# server
rm -rf quizpb/*.go
protoc --go_out=quizpb --go_opt=paths=source_relative \
    --go-grpc_out=quizpb --go-grpc_opt=paths=source_relative \
    --proto_path=quizpb \
    quizpb/*.proto

