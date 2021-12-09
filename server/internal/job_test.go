package internal

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Hexcles/Vaporeon/protos"
	"github.com/Hexcles/Vaporeon/worker"
)

func TestJobToPb_running(t *testing.T) {
	id := &pb.JobId{Uuid: uuid.NewString()}
	owner := "guest"
	args := []string{"echo", "hello"}
	now := time.Now()
	job := &Job{
		Owner: owner,
		Job: &worker.Job{
			Args:    args,
			Started: now,
		},
	}
	got := jobToPb(id, job)
	want := &pb.Job{
		Owner:   owner,
		Id:      id,
		Args:    args,
		Started: timestamppb.New(now),
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("jobToPb() returned unexpected proto; diff (-want +got):\n%s", diff)
	}
}

func TestJobToPb_stopped(t *testing.T) {
	id := &pb.JobId{Uuid: uuid.NewString()}
	j, err := worker.Launch([]string{"sleep", "100000"})
	if err != nil {
		t.Fatalf("Launch() returned unexpected error: %v", err)
	}
	job := &Job{Owner: "guest", Job: j}
	if err := j.Kill(); err != nil {
		t.Fatalf("Kill() returned unexpected error: %v", err)
	}
	got := jobToPb(id, job)
	if got.ExitCode == 0 {
		t.Error("ExitCode is 0; want non-zero")
	}
	if got.Stopped.AsTime().IsZero() {
		t.Error("Stopped is zero; want non-zero")
	}
}
