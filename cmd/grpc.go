package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/mhasnanr/ewallet-ums/bootstrap"
	pb "github.com/mhasnanr/ewallet-ums/cmd/tokenvalidation"
	"github.com/mhasnanr/ewallet-ums/helpers"
	handler "github.com/mhasnanr/ewallet-ums/internal/handler/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ServeGRPC() {
	grpcPort := bootstrap.GetEnv("GRPC_PORT", "7000")
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("failed to listen grpc port: ", err)
	}

	server := grpc.NewServer()

	jwtManager := &helpers.JWTManager{}

	tokenHandler := handler.NewTokenValidationHandler(jwtManager)
	pb.RegisterTokenValidationServer(server, tokenHandler)

	reflection.Register(server)

	fmt.Printf("gRPC server is running on port %s...\n", grpcPort)
	if err := server.Serve(listener); err != nil {
		log.Fatal("failed to serve grpc port: ", err)
	}
}
