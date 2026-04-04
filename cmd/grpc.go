package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/mhasnanr/e-wallet/bootstrap"
	"google.golang.org/grpc"
)

func ServeGRPC() {
	listener, err := net.Listen("tcp", ":"+bootstrap.GetEnv("GRPC_PORT", "7000"))
	if err != nil {
		log.Fatal("failed to listen grpc port: ", err)
	}

	server := grpc.NewServer()

	if err := server.Serve(listener); err != nil {
		fmt.Println("pe")
		log.Fatal("failed to serve grpc port: ", err)
	}
}
