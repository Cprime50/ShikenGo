TODO:
[] add password  when generating certificate.pem

[] write tests

[x] generate proto

[] switch make file to main folder and finish the script


WHne installing proto-gen-go-grpc on linux for the first time

go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

go get google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go

PATH="${PATH}:${HOME}/go/bin"
