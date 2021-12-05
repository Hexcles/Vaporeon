package worker

import (
	"io/ioutil"
	"os"
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
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			job, err := Launch(c.args)
			if err != nil {
				t.Fatalf("Launch() returned non-nil error: %v", err)
			}

			if got, err := ioutil.ReadAll(job.Stdout()); string(got) != c.wantOut || err != nil {
				t.Errorf("ReadAll(Stdout()) = %q, %v; want %q, nil", string(got), err, c.wantOut)
			}
			if got, err := ioutil.ReadAll(job.Stderr()); string(got) != c.wantErr || err != nil {
				t.Errorf("ReadAll(Stderr()) = %q, %v; want %q, nil", string(got), err, c.wantErr)
			}
		})
	}
}

func TestLaunch_Error(t *testing.T) {
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

func TestKill(t *testing.T) {
	job, err := Launch([]string{"sleep", "1000000"})
	if err != nil {
		t.Fatalf("Launch() returned non-nil error: %v", err)
	}
	if err := job.Kill(); err != nil {
		t.Errorf("Kill() returned non-nil error: %v", err)
	}
	// This should not block.
	ioutil.ReadAll(job.Stdout())
}

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}
