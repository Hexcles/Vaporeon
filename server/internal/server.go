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

// Auther is a simple authorization interface.
type Auther interface {
	GetPeerID(context.Context) (string, error)
	CanManage(context.Context, string) (bool, error)
	CanShutdown(context.Context) (bool, error)
}

// Server implements JobWorkerServer.
type Server struct {
	auther    Auther
	jobs      sync.Map
	shutdown  chan<- struct{}
	closeOnce sync.Once

	pb.UnimplementedJobWorkerServer
}

// New creates a new server.
func New(auther Auther, shutdown chan<- struct{}) *Server {
	return &Server{auther: auther, shutdown: shutdown}
}

// Launch launches a new job.
func (s *Server) Launch(ctx context.Context, req *pb.Job) (*pb.JobId, error) {
	email, err := s.auther.GetPeerID(ctx)
	if err != nil {
		return nil, err
	}
	job, err := worker.Launch(req.Args)
	if err != nil {
		log.Printf("Failed to launched job %v: %s", req.Args, err)
		return nil, err
	}
	jid := &pb.JobId{Uuid: uuid.NewString()}
	s.saveJob(jid, &Job{Job: job, Owner: email})
	log.Printf("Launched job %s for %s: %v", jid, email, req.Args)
	return jid, nil
}

// Kill sends SIGKILL to the job and blocks until the job exits.
func (s *Server) Kill(ctx context.Context, req *pb.JobId) (*pb.Job, error) {
	job, err := s.loadJob(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := job.Job.Kill(); err == nil {
		log.Printf("Killed job %s %v", req.Uuid, job.Job.Args)
	} else {
		log.Printf("Failed to kill job %s %v: %s", req.Uuid, job.Job.Args, err)
	}
	return jobToPb(req, job), err
}

// Query returns the status of the job.
func (s *Server) Query(ctx context.Context, req *pb.JobId) (*pb.Job, error) {
	job, err := s.loadJob(ctx, req)
	if err != nil {
		return nil, err
	}
	return jobToPb(req, job), nil
}

// StreamOutput streams the output for the job.
func (s *Server) StreamOutput(req *pb.JobId, resp pb.JobWorker_StreamOutputServer) error {
	job, err := s.loadJob(resp.Context(), req)
	if err != nil {
		return err
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
func (s *Server) Shutdown(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	ok, err := s.auther.CanShutdown(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "no permission")
	}
	s.closeOnce.Do(s.close)
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

func (s *Server) saveJob(id *pb.JobId, job *Job) {
	// Uuid.NewString should always succeed.
	if id.Uuid == "" {
		panic("Got empty UUID")
	}
	s.jobs.Store(id.Uuid, job)
}

func (s *Server) loadJob(ctx context.Context, id *pb.JobId) (*Job, error) {
	job, ok := s.jobs.Load(id.Uuid)
	if !ok {
		return nil, status.Error(codes.NotFound, "job not found or no permission")
	}
	j := job.(*Job)
	ok, err := s.auther.CanManage(ctx, j.Owner)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.NotFound, "job not found or no permission")
	}
	return j, nil
}

// close closes the shutdown channel. This is behind a sync.Once to avoid
// panics from closing the channel multiple times .
func (s *Server) close() {
	close(s.shutdown)
}
