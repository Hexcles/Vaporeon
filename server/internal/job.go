package internal

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Hexcles/Vaporeon/protos"
	"github.com/Hexcles/Vaporeon/worker"
)

// Job represents an owned job.
type Job struct {
	Owner string
	Job   *worker.Job
}

func jobToPb(id *pb.JobId, job *Job) *pb.Job {
	ret := &pb.Job{
		Owner:   job.Owner,
		Id:      id,
		Args:    job.Job.Args,
		Started: timestamppb.New(job.Job.Started),
	}
	if s := job.Job.Status(); !s.Stopped.IsZero() {
		ret.Stopped = timestamppb.New(s.Stopped)
		if s.ExitSignal != 0 {
			ret.ExitSignal = int32(s.ExitSignal)
		} else {
			ret.ExitCode = int32(s.ExitCode)
		}
	}
	return ret
}
