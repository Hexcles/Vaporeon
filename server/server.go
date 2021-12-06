package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/Hexcles/Vaporeon/protos"
	"github.com/Hexcles/Vaporeon/server/internal"
	"github.com/Hexcles/Vaporeon/worker"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 50051, "The server port")
)

func init() {
	worker.Init()
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	shutdown := make(chan struct{})
	server := internal.New(shutdown)
	pb.RegisterJobWorkerServer(grpcServer, server)
	go func() {
		<-shutdown
		log.Println("Shutting down")
		// Do not attempt to GracefulStop here as open streams may
		// block until the launched processes exit.
		grpcServer.Stop()
		server.KillAll()
	}()
	grpcServer.Serve(lis)
}
