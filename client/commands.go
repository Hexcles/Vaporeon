package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Hexcles/Vaporeon/protos"
)

// Launch launches a job on the server.
func Launch(ctx context.Context, client pb.JobWorkerClient, args []string) (*pb.JobId, error) {
	log.Printf("Launching %v", args)
	id, err := client.Launch(ctx, &pb.Job{Args: args})
	if err != nil {
		return nil, err
	}
	log.Printf("Launched job %s", id.Uuid)
	return id, nil
}

// Query prints the status of a job.
func Query(ctx context.Context, client pb.JobWorkerClient, id string) error {
	log.Printf("Querying %s", id)
	job, err := client.Query(ctx, &pb.JobId{Uuid: id})
	if err != nil {
		return err
	}
	fmt.Println(job)
	return nil
}

// Stream prints the stdout and stderr of a job to stdout and stderr
// respectively. It returns when the job finishes.
func Stream(ctx context.Context, client pb.JobWorkerClient, id string) error {
	log.Printf("Streaming output from job %s", id)
	stream, err := client.StreamOutput(ctx, &pb.JobId{Uuid: id})
	if err != nil {
		return err
	}
	for {
		output, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		if output.Stdout != nil {
			os.Stdout.Write(output.Stdout)
		}
		if output.Stderr != nil {
			os.Stderr.Write(output.Stderr)
		}
	}
}

// Kill requests the server to kill a job.
func Kill(ctx context.Context, client pb.JobWorkerClient, id string) error {
	log.Printf("Killing %s", id)
	job, err := client.Kill(ctx, &pb.JobId{Uuid: id})
	if err != nil {
		return err
	}
	log.Printf("Killed job: %s", job)
	return nil
}

// Shutdown closes the server and kills all jobs.
func Shutdown(ctx context.Context, client pb.JobWorkerClient) error {
	log.Println("Shutting down the server")
	_, err := client.Shutdown(ctx, &emptypb.Empty{})
	return err
}
