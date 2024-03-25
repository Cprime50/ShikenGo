[] find what password it is requesting for when generating certificate.pem

[] write tests

[x] generate proto

[] switch make file to main folder and finish the script

[] Look at how ssg did his db and service stuff, see the way he used go routine to query faster, use his style of db and service sepeartion, refactor code to look like that

WHne installing proto-gen-go-grpc on linux for the first time

go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

go get google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go

PATH="${PATH}:${HOME}/go/bin"

This chathtp link https://chat.openai.com/c/3236d1ed-7d86-42d8-90d1-a6b615c4239c to finish the quiz section