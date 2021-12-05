// Package cgroup is a simple, incomplete encapsulation of cgroupv2.
//
// Users of this package are expected to launch the process in a cgroup it can
// control: the process needs to be able to write to the cgroup hierarchy and
// is the only process in this hierarchy. This can be achieved with
// systemd-run --scope.
//
// It is hardcoded to use memory, io and cpu controllers. If any of these
// controllers are unavailable, the error is intentionally ignored (i.e.
// best-effort support for individual controllers).
//
// This package should be used as early as possible, before any potential
// subprocess is created (e.g. during the init() of the main package).
package cgroup

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var wantControllers = []string{"memory", "io", "cpu"}

const (
	// maxMemory is written to memory.max.
	maxMemory = 128 * 1024 * 1024 // 128MiB
	// maxIOBps is written to io.max as rbps and wbps for each block device.
	maxIOBps = 2 * 1024 * 1024 // 2MiB/s
	// maxCPU is written to cpu.max.
	maxCPU = "20000 100000" // 20%
)

// TODO: this is the default of systemd, but not guaranteed on all systems.
const cgroupRoot = "/sys/fs/cgroup"

// Cgroup remembers a cgroup hierarchy, which won't change even if the process
// is moved to a new hierarchy.
type Cgroup struct {
	Path string
	leaf bool
}

// New creates a new Cgroup at the current hierarchy of the process. It returns
// an error if the current cgroup hierarchy cannot be found.
func New() (*Cgroup, error) {
	path, err := getCgroupPath("/proc/self/cgroup")
	if err != nil {
		return nil, err
	}
	return &Cgroup{path, true}, nil
}

// Check returns an error if the cgroup hierarchy cannot be managed by this
// process (e.g. if there are other processes in the same hierarchy), in which
// case cgroup setup should be aborted.
//
// It also checks available controllers and will print a log message if a
// desired controller is unavailable (but do not return an error).
func (c *Cgroup) Check() error {
	procs, err := os.ReadFile(filepath.Join(c.Path, "cgroup.procs"))
	if err != nil {
		return fmt.Errorf("cgroup: failed to read cgroup.procs: %s", err)
	}
	if strings.Count(string(procs), "\n") != 1 {
		return errors.New("cgroup: cannot control the cgroup hierarchy; found other processes")
	}
	content, err := os.ReadFile(filepath.Join(c.Path, "cgroup.controllers"))
	if err != nil {
		return fmt.Errorf("cgroup: failed to read cgroup.controllers: %s", err)
	}
	cMap := make(map[string]bool)
	for _, c := range strings.Split(strings.TrimSpace(string(content)), " ") {
		cMap[c] = true
	}
	for _, c := range wantControllers {
		if !cMap[c] {
			log.Printf("cgroup: %s controller is not available and will be skipped", c)
		}
	}
	return nil
}

// MoveToNewSubtree creates a new cgroup hierarchy and moves the current
// process there. The given subtree is relative to the *saved* hierarchy (not
// the current, if the process is already moved).
func (c *Cgroup) MoveToNewSubtree(subtree string) (*Cgroup, error) {
	newPath := filepath.Join(c.Path, subtree)
	if err := syscall.Mkdir(newPath, 0700); err != nil {
		return nil, fmt.Errorf("cgroup: failed to create %s: %s", subtree, err)
	}
	if err := os.WriteFile(
		filepath.Join(newPath, "cgroup.procs"),
		[]byte(fmt.Sprintf("%d", os.Getpid())),
		0600); err != nil {
		return nil, fmt.Errorf("cgroup: failed to move into %s: %s", subtree, err)
	}
	c.leaf = false
	return &Cgroup{newPath, true}, nil
}

// EnableSubtreeControl enables subtree control (i.e. delegates all wanted
// controllers to the subtree). Only call this method after a successful call
// to MoveToNewSubtree since the kernel doesn't allow internal/non-leaf
// processes in cgroupv2.
//
// Errors are intentionally ignored since they imply unavailability of
// controllers.
func (c *Cgroup) EnableSubtreeControl() {
	if c.leaf {
		panic("EnableSubtreeControl called before MoveToNewSubtree")
	}
	// Enable controllers one by one so an error won't affect others.
	for _, controller := range wantControllers {
		_ = os.WriteFile(
			filepath.Join(c.Path, "cgroup.subtree_control"),
			[]byte("+"+controller),
			0600)
	}
}

// EnableLimits enables the hard-coded limits on all wanted controllers.
//
// Errors are intentionally ignored since they imply unavailability of
// controllers.
func (c *Cgroup) EnableLimits() {
	// Memory
	_ = os.WriteFile(
		filepath.Join(c.Path, "memory.max"),
		[]byte(fmt.Sprintf("%d", maxMemory)),
		0600)
	// CPU
	_ = os.WriteFile(
		filepath.Join(c.Path, "cpu.max"),
		[]byte(maxCPU),
		0600)
	// IO
	for _, maj := range findBlockMajors("/proc/partitions") {
		_ = os.WriteFile(
			filepath.Join(c.Path, "io.max"),
			[]byte(fmt.Sprintf("%s:0 rbps=%d wbps=%d", maj, maxIOBps, maxIOBps)),
			0600)
	}
}

func getCgroupPath(path string) (string, error) {
	read, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cgroup: failed to read /proc/self/cgroup: %s", err)
	}
	cgroup := strings.TrimSpace(string(read))
	// hierarchy-ID:controller-list:cgroup-path
	split := strings.Split(cgroup, ":")
	if len(split) != 3 {
		return "", fmt.Errorf("cgroup: unable to parse %q", cgroup)
	}
	return filepath.Join(cgroupRoot, split[2]), nil
}

func findBlockMajors(partitions string) []string {
	content, err := os.ReadFile(partitions)
	if err != nil {
		log.Printf("cgroup: failed to read %s: %s", partitions, err)
		return nil
	}
	lines := strings.Split(string(content), "\n")
	// The first line is the header.
	if len(lines) == 0 {
		return nil
	}
	majors := make(map[string]bool)
	for _, line := range lines[1:] {
		split := strings.Fields(line)
		if len(split) == 4 {
			majors[split[0]] = true
		}
	}
	ret := make([]string, 0, len(majors))
	for key := range majors {
		ret = append(ret, key)
	}
	return ret
}
