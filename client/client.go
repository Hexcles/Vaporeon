package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/Hexcles/Vaporeon/protos"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	flag.Parse()
	args := flag.Args()

	var opts []grpc.DialOption
	if *tls {
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewJobWorkerClient(conn)
	ctx := context.Background()

	if len(args) == 0 {
		id, err := Launch(ctx, client, []string{"ls", "/"})
		if err != nil {
			log.Fatal(err)
		}
		Stream(ctx, client, id.Uuid)
		Query(ctx, client, id.Uuid)
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
		err = Kill(ctx, client, args[1])
	case "shutdown":
		err = Shutdown(ctx, client)
	default:
		log.Fatalf("Unknown command: %s", args[0])
	}
	if err != nil {
		log.Println(err)
	}
}
