package worker

import (
	"io"
	"os"
	"syscall"
	"testing"
)

func TestLaunch(t *testing.T) {
	type testCase struct {
		name    string
		args    []string
		wantOut string
		wantErr string
	}
	cases := []testCase{
		{name: "basic IO", args: []string{"echo", "-n", "hello"}, wantOut: "hello"},
		{name: "UID mapping", args: []string{"id", "-u"}, wantOut: "0\n"},
		{name: "GID mapping", args: []string{"id", "-g"}, wantOut: "0\n"},
		{name: "hostname", args: []string{"hostname"}, wantOut: hostname + "\n"},
	}
	for _, c := range cases {
		// Run the cases in parallel to prove goroutine-safety. Since
		// they run in goroutines, capture the loop variable here.
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			job, err := Launch(c.args)
			if err != nil {
				t.Fatalf("Launch() returned non-nil error: %v", err)
			}
			if got, err := io.ReadAll(job.Stdout()); string(got) != c.wantOut || err != nil {
				t.Errorf("ReadAll(Stdout()) = %q, %v; want %q, nil", string(got), err, c.wantOut)
			}
			if got, err := io.ReadAll(job.Stderr()); string(got) != c.wantErr || err != nil {
				t.Errorf("ReadAll(Stderr()) = %q, %v; want %q, nil", string(got), err, c.wantErr)
			}
			if s := job.Status(); s.Stopped.IsZero() || s.ExitSignal != 0 || s.ExitCode != 0 {
				t.Errorf("Status() = %+v; want {Stopped:non-zero, ExitSignal:0, ExitCode:0}", s)
			}
		})
	}
}

func TestLaunch_error(t *testing.T) {
	type testCase struct {
		name string
		args []string
	}
	cases := []testCase{
		{name: "empty args"},
		{name: "empty args[0]", args: []string{""}},
		{name: "executable not found", args: []string{"tHiS_ShOuLd_NoT_ExIsT"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Launch(c.args)
			if err == nil {
				t.Errorf("Launch(%v) returned nil error", c.args)
			}
		})
	}
}

func TestLaunch_exit_with_nonzero(t *testing.T) {
	job, err := Launch([]string{"false"})
	if err != nil {
		t.Errorf("Launch() returned non-nil error: %v", err)
	}
	_, _ = io.ReadAll(job.Stderr())
	if s := job.Status(); s.Stopped.IsZero() || s.ExitSignal != 0 || s.ExitCode != 1 {
		t.Errorf("Status() = %+v; want {Stopped:non-zero, ExitSignal:0, ExitCode:1}", s)
	}
}

func TestKill(t *testing.T) {
	job, err := Launch([]string{"sleep", "1000000"})
	if err != nil {
		t.Fatalf("Launch() returned non-nil error: %v", err)
	}
	if s := job.Status(); !s.Stopped.IsZero() {
		t.Error("Before Kill(), Status() returned non-zero Stopped")
	}
	if err := job.Kill(); err != nil {
		t.Errorf("Kill() returned non-nil error: %v", err)
	}
	// This should not block.
	if _, err = io.ReadAll(job.Stdout()); err != nil {
		t.Fatalf("ReadAll(Stdout()) returned non-nil error: %v", err)
	}
	if s := job.Status(); s.Stopped.IsZero() || s.ExitSignal != syscall.SIGKILL || s.ExitCode != 0 {
		t.Errorf("After Kill(), Status() = %+v; want {Stopped:non-zero, ExitSignal:SIGKILL, ExitCode:0}", s)
	}
}

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
