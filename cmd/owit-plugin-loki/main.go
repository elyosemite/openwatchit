package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	apiv1 "github.com/elyosemite/openwatchit/api/v1"
	"github.com/elyosemite/openwatchit/plugins/loki"
)

func main() {
	addr := "127.0.0.1:50051"
	if v := os.Getenv("OWIT_PLUGIN_ADDR"); v != "" {
		addr = v
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen error: %v\n", err)
		os.Exit(1)
	}

	srv := grpc.NewServer()
	apiv1.RegisterPluginServiceServer(srv, loki.NewServer())

	fmt.Printf("owit-plugin-loki listening on %s\n", addr)
	if err := srv.Serve(lis); err != nil {
		fmt.Fprintf(os.Stderr, "serve error: %v\n", err)
		os.Exit(1)
	}
}
