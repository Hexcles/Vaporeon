package internal

import (
	"io"
	"syscall"
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

func TestJobToPb_killed(t *testing.T) {
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
	if got.ExitSignal != int32(syscall.SIGKILL) {
		t.Errorf("ExitSignal is %d; want %d", got.ExitSignal, syscall.SIGKILL)
	}
	if got.ExitCode != 0 {
		t.Errorf("ExitCode is %d; want 0", got.ExitCode)
	}
	if got.Stopped.AsTime().IsZero() {
		t.Error("Stopped is zero; want non-zero")
	}
}

func TestJobToPb_exited(t *testing.T) {
	id := &pb.JobId{Uuid: uuid.NewString()}
	j, err := worker.Launch([]string{"false"})
	if err != nil {
		t.Fatalf("Launch() returned unexpected error: %v", err)
	}
	job := &Job{Owner: "guest", Job: j}
	_, _ = io.ReadAll(j.Stderr())
	got := jobToPb(id, job)
	if got.ExitSignal != 0 {
		t.Errorf("ExitSignal is %d; want 0", got.ExitSignal)
	}
	if got.ExitCode != 1 {
		t.Errorf("ExitCode is %d; want 1", got.ExitCode)
	}
	if got.Stopped.AsTime().IsZero() {
		t.Error("Stopped is zero; want non-zero")
	}
}
