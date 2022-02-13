package internal

import (
	"context"
	"testing"

	"go.uber.org/goleak"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Hexcles/Vaporeon/protos"
	"github.com/Hexcles/Vaporeon/worker"
	"github.com/google/uuid"
)

type fakeAuther struct {
	id          string
	canManage   bool
	canShutdown bool
}

func (f fakeAuther) GetPeerID(context.Context) (string, error) {
	return f.id, nil
}

func (f fakeAuther) CanManage(context.Context, string) (bool, error) {
	return f.canManage, nil
}

func (f fakeAuther) CanShutdown(context.Context) (bool, error) {
	return f.canShutdown, nil
}

func TestLaunch(t *testing.T) {
	ctx := context.Background()
	owner := "guest"
	s := New(fakeAuther{id: owner}, make(chan struct{}))
	if _, err := s.Launch(ctx, &pb.Job{Args: []string{"echo", "hello"}}); err != nil {
		t.Fatalf("Launch() returned non-nil error: %v", err)
	}
	// The job will be waited in a goroutine, upon which all server
	// resources will be garbage collected.
}

func TestQuery(t *testing.T) {
	ctx := context.Background()
	id := &pb.JobId{Uuid: uuid.NewString()}
	auther := &fakeAuther{canManage: true}
	s := New(auther, make(chan struct{}))
	s.saveJob(id, &Job{Owner: "guest", Job: &worker.Job{}})

	// Success:
	if _, err := s.Query(ctx, id); err != nil {
		t.Errorf("Query() returned non-nil error: %v", err)
	}

	// Not found:
	_, err := s.Query(ctx, &pb.JobId{Uuid: uuid.NewString()})
	st, _ := status.FromError(err)
	if st.Code() != codes.NotFound {
		t.Errorf("Query(non-existent ID) returned error code %s; want %s", st.Code(), codes.NotFound)
	}

	// Permission denied:
	auther.canManage = false
	_, err = s.Query(ctx, id)
	st, _ = status.FromError(err)
	if st.Code() != codes.NotFound {
		t.Errorf("Query(no permission) returned error code %s; want %s", st.Code(), codes.NotFound)
	}
}

func TestKill(t *testing.T) {
	ctx := context.Background()
	auther := &fakeAuther{id: "guest", canManage: true}
	s := New(auther, make(chan struct{}))
	id, err := s.Launch(ctx, &pb.Job{Args: []string{"sleep", "100000"}})
	if err != nil {
		t.Fatalf("Launch() returned non-nil error: %v", err)
	}

	// Success:
	if _, err := s.Kill(ctx, id); err != nil {
		t.Errorf("Kill() returned non-nil error: %v", err)
	}

	// Not found:
	_, err = s.Kill(ctx, &pb.JobId{Uuid: uuid.NewString()})
	st, _ := status.FromError(err)
	if st.Code() != codes.NotFound {
		t.Errorf("Kill(non-existent ID) returned error code %s; want %s", st.Code(), codes.NotFound)
	}

	// Permission denied:
	auther.canManage = false
	_, err = s.Kill(ctx, id)
	st, _ = status.FromError(err)
	if st.Code() != codes.NotFound {
		t.Errorf("Kill(no permission) returned error code %s; want %s", st.Code(), codes.NotFound)
	}
}

func TestShutdown_success(t *testing.T) {
	ch := make(chan struct{}, 1)
	s := New(fakeAuther{canShutdown: true}, ch)
	if _, err := s.Shutdown(context.Background(), nil); err != nil {
		t.Fatalf("Shutdown() returned non-nil error: %v", err)
	}
	<-ch
}

func TestShutdown_permission_denied(t *testing.T) {
	ch := make(chan struct{}, 1)
	s := New(fakeAuther{canManage: true}, ch)
	_, err := s.Shutdown(context.Background(), nil)
	st, _ := status.FromError(err)
	if st.Code() != codes.PermissionDenied {
		t.Fatalf("Shutdown() returned error code %s; want %s", st.Code(), codes.PermissionDenied)
	}
	select {
	case <-ch:
		t.Error("shutdown channel unexpectedly returned")
	default:
		// Do not block.
	}
}

func TestShutdown_multiple(t *testing.T) {
	// Calling Shutdown() multiple times should not panic.
	ch := make(chan struct{}, 1)
	s := New(fakeAuther{canShutdown: true}, ch)
	if _, err := s.Shutdown(context.Background(), nil); err != nil {
		t.Fatalf("Shutdown() returned non-nil error: %v", err)
	}
	if _, err := s.Shutdown(context.Background(), nil); err != nil {
		t.Fatalf("Shutdown() returned non-nil error: %v", err)
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
