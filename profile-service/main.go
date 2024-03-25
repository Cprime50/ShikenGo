package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/Cprime50/user/db"
	pb "github.com/Cprime50/user/profilepb"
	"github.com/Cprime50/user/src"
	"github.com/Cprime50/user/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	//"google.golang.org/grpc/reflection"
	//"google.golang.org/grpc/reflection"
)

// func newServer() *server {
// 	// Initialize and return a new instance of your profile server
// 	return &server{}
// }

var (
	_         = utils.LoadEnv()
	ENV       = utils.MustHaveEnv("ENV")
	GRPC_PORT = utils.MustHaveEnv("GRPC_PORT")
	CERT_PATH = utils.MustHaveEnv("CERT_PATH")
	KEY_PATH  = utils.MustHaveEnv("KEY_PATH")
)

func main() {
	server := &src.Server{}

	//Load logger

	//Connect db
	Db, err := db.Connect()
	if err != nil {
		slog.Error("Error opening database", "db.Connect", err)
		log.Fatal("Error connecting to Db", err)
	}
	log.Println("Database connected successfully")

	// migrations
	log.Printf("Migrations Started")
	err = db.Migrate(Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Db.Close()

	// Run the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%v", GRPC_PORT))
	if err != nil {
		slog.Error("Error listening on gRPC port", "net.Listen", err)
		panic(err)
	}

	var s *grpc.Server
	if ENV == "production" {
		certificate, err := tls.LoadX509KeyPair(CERT_PATH, KEY_PATH)
		if err != nil {
			slog.Error("Error loading TLS certificate", "tls.LoadX509KeyPair", err)
			panic(err)
		}
		s = grpc.NewServer(grpc.Creds(credentials.NewServerTLSFromCert(&certificate)))
	} else {
		//grpcLogger := grpc.UnaryInterceptor(utils.GrpcLogger)
		s = grpc.NewServer()
	}

	//}
	reflection.Register(s)
	pb.RegisterProfileServiceServer(s, server)

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
		panic(err)
	}

}
