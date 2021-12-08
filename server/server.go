package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Hexcles/Vaporeon/certs"
	pb "github.com/Hexcles/Vaporeon/protos"
	"github.com/Hexcles/Vaporeon/server/internal"
	"github.com/Hexcles/Vaporeon/worker"
)

var (
	insecure = flag.Bool("insecure", false, "Disable TLS")
	caFile   = flag.String("ca_file", "../certs/ca_cert.pem", "The file containing the CA root cert file")
	certFile = flag.String("cert_file", "../certs/server_cert.pem", "The TLS cert file")
	keyFile  = flag.String("key_file", "../certs/server_key.pem", "The TLS key file")
	port     = flag.Int("port", 50051, "The server port")
)

func init() {
	worker.Init()
}

func loadKeyPair() (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		return nil, err
	}
	capool, err := certs.LoadCA(*caFile)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig), nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if !*insecure {
		creds, err := loadKeyPair()
		if err != nil {
			log.Fatalf("Failed to load TLS credentials: %s", err)
		}
		opts = append(opts, grpc.Creds(creds))
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
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
