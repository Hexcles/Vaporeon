package internal

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/Hexcles/Vaporeon/protos"
	"github.com/Hexcles/Vaporeon/worker"
)

// Server implements JobWorkerServer.
type Server struct {
	jobs     sync.Map
	shutdown chan<- struct{}

	pb.UnimplementedJobWorkerServer
}

// New creates a new server.
func New(shutdown chan<- struct{}) *Server {
	return &Server{shutdown: shutdown}
}

// Launch launches a new job.
func (s *Server) Launch(ctx context.Context, req *pb.Job) (*pb.JobId, error) {
	id := uuid.NewString()
	job, err := worker.Launch(req.Args)
	if err != nil {
		log.Printf("Failed to launched job %v: %s", req.Args, err)
		return nil, err
	}
	s.jobs.Store(id, &Job{Job: job})
	log.Printf("Launched job %s: %v", id, req.Args)
	return &pb.JobId{Uuid: id}, nil
}

// Kill sends SIGKILL to the job and blocks until the job exits.
func (s *Server) Kill(ctx context.Context, req *pb.JobId) (*pb.Job, error) {
	job := s.loadJob(req)
	if job == nil {
		return nil, status.Error(codes.NotFound, "job not found")
	}
	err := job.Job.Kill()
	if err == nil {
		log.Printf("Killed job %s %v", req.Uuid, job.Job.Args)
	} else {
		log.Printf("Failed to kill job %s %v: %s", req.Uuid, job.Job.Args, err)
	}
	return jobToPb(req, job), err
}

// Query returns the status of the job.
func (s *Server) Query(ctx context.Context, req *pb.JobId) (*pb.Job, error) {
	job := s.loadJob(req)
	if job == nil {
		return nil, status.Error(codes.NotFound, "job not found")
	}
	return jobToPb(req, job), nil
}

// StreamOutput streams the output for the job.
func (s *Server) StreamOutput(req *pb.JobId, resp pb.JobWorker_StreamOutputServer) error {
	job := s.loadJob(req)
	if job == nil {
		return status.Error(codes.NotFound, "job not found")
	}
	errCh1 := make(chan error, 1)
	errCh2 := make(chan error, 1)
	go sendBuffer(resp.Context(), job.Job.Stdout(), errCh1, func(b []byte) error {
		return resp.Send(&pb.Output{Stdout: b})
	})
	go sendBuffer(resp.Context(), job.Job.Stderr(), errCh2, func(b []byte) error {
		return resp.Send(&pb.Output{Stderr: b})
	})
	return waitForSenders(resp.Context(), errCh1, errCh2)
}

// Shutdown signals to shut down the server and kill all jobs.
func (s *Server) Shutdown(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.shutdown <- struct{}{}
	return &emptypb.Empty{}, nil
}

// KillAll kills all jobs.
//
// Make sure to stop accepting new jobs before calling this method.
// Errors are silently ignored.
func (s *Server) KillAll() {
	s.jobs.Range(func(key, value interface{}) bool {
		job := value.(*Job)
		// Some jobs may have ended.
		_ = job.Job.Kill()
		return true
	})
}

func (s *Server) loadJob(id *pb.JobId) *Job {
	job, ok := s.jobs.Load(id.Uuid)
	if !ok {
		return nil
	}
	return job.(*Job)
}
