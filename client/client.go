package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Hexcles/Vaporeon/certs"
	pb "github.com/Hexcles/Vaporeon/protos"
)

var (
	insecure   = flag.Bool("insecure", false, "Disable TLS")
	caFile     = flag.String("ca_file", "../certs/ca_cert.pem", "The file containing the CA root cert file")
	certFile   = flag.String("cert_file", "../certs/client1_cert.pem", "The TLS cert file")
	keyFile    = flag.String("key_file", "../certs/client1_key.pem", "The TLS key file")
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `Commands:
  launch [command args...]
  query [ID]
  kill [ID]
  stream [ID]
  shutdown
`)
	}
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
		Certificates: []tls.Certificate{certificate},
		RootCAs:      capool,
	}
	return credentials.NewTLS(tlsConfig), nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	var opts []grpc.DialOption
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		creds, err := loadKeyPair()
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer conn.Close()
	client := pb.NewJobWorkerClient(conn)
	ctx := context.Background()

	if len(args) == 0 {
		id, err := Launch(ctx, client, []string{"ls", "/"})
		if err != nil {
			log.Fatal(err)
		}
		if err := Stream(ctx, client, id.Uuid); err != nil {
			log.Fatal(err)
		}
		if err := Query(ctx, client, id.Uuid); err != nil {
			log.Fatal(err)
		}
		return
	}

	if args[0] != "shutdown" && len(args) == 1 {
		log.Fatal("Insufficient arguments")
	}
	switch args[0] {
	case "launch":
		_, err = Launch(ctx, client, args[1:])
	case "query":
		err = Query(ctx, client, args[1])
	case "kill":
		err = Kill(ctx, client, args[1])
	case "stream":
		err = Stream(ctx, client, args[1])
	case "shutdown":
		err = Shutdown(ctx, client)
	default:
		log.Fatalf("Unknown command: %s", args[0])
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
