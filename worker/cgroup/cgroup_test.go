package cgroup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestGetCgroupPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cgroup")
	err := os.WriteFile(path, []byte("0::/user.slice/user-1000.slice\n"), 0600)
	if err != nil {
		t.Fatal(err)
	}
	want := "/sys/fs/cgroup/user.slice/user-1000.slice"
	cgroup, err := getCgroupPath(path)
	if cgroup != want || err != nil {
		t.Errorf("getCgroupPath() = %q, %v; want %q, nil", cgroup, err, want)
	}
}

func TestFindBlockMajors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "partitions")
	err := os.WriteFile(path, []byte(`major minor  #blocks  name

259        0  500107608 nvme0n1
259        1     266240 nvme0n1p1
254        0  341187571 dm-0
`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	majors := findBlockMajors(path)
	sort.Strings(majors)
	if len(majors) != 2 || majors[0] != "254" || majors[1] != "259" {
		t.Errorf("findBlockMajors() = %v; want [254, 259]", majors)
	}
}

func TestCgroupCheck_invalidDir(t *testing.T) {
	fake := t.TempDir()
	cgroup := &Cgroup{fake, true}
	if err := cgroup.Check(); err == nil {
		t.Error("Check() did not return an error")
	}
}

func TestCgroupCheck_manyProcs(t *testing.T) {
	fake := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(fake, "cgroup.procs"),
		[]byte("1000\n1001\n"),
		0600); err != nil {
		t.Fatal(err)
	}
	cgroup := &Cgroup{fake, true}
	if err := cgroup.Check(); err == nil {
		t.Error("Check() did not return an error")
	}
}

func TestCgroupCheck_success(t *testing.T) {
	fake := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(fake, "cgroup.procs"),
		[]byte("1000\n"),
		0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(fake, "cgroup.controllers"),
		[]byte("memory io cpu\n"),
		0600); err != nil {
		t.Fatal(err)
	}
	cgroup := &Cgroup{fake, true}
	if err := cgroup.Check(); err != nil {
		t.Errorf("Check() = %v; want nil", err)
	}
}

func TestMoveToNewSubtree(t *testing.T) {
	fake := t.TempDir()
	cgroup := &Cgroup{fake, true}
	cgroup.MoveToNewSubtree("new")
	read, err := os.ReadFile(filepath.Join(fake, "new", "cgroup.procs"))
	if err != nil {
		t.Fatalf("No new cgroup.procs: %s", err)
	}
	want := fmt.Sprintf("%d", os.Getpid())
	if string(read) != want {
		t.Errorf("cgroup.procs = %s; want %s", string(read), want)
	}
}
