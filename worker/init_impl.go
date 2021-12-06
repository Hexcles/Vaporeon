package worker

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"

	"github.com/Hexcles/Vaporeon/worker/cgroup"
)

const specialArg0 = "<init>"
const hostname = "vaporeon"
const cgroupName = "_vaporeon"

func fatal(err error) {
	log.Fatalf("Init error: %s", err)
}

func parentInit() {
	// This is for extra safety in case of thread-mode cgroup.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := setUpParentCgroup(); err != nil {
		log.Printf("Init: skipping cgroup setup: %s", err)
	}
}

// setUpParentCgroup produces a non-fatal error if cgroup setup fails.
//
// If successful, we get the following hierarchy:
//   - original
//     subtree_control: on
//     - cgroupName
//       procs: [PID]
func setUpParentCgroup() error {
	c, err := cgroup.New()
	if err != nil {
		return err
	}
	if err := c.Check(); err != nil {
		return err
	}
	if _, err := c.MoveToNewSubtree(cgroupName); err != nil {
		return err
	}
	c.EnableSubtreeControl()
	return nil
}

func childInit() {
	// We are going to call unshare which works on the thread level.
	runtime.LockOSThread()
	setUpChildCgroup()
	// Even if cgroup setup is skipped, we still want to create a new
	// namespace for isolation.
	if err := syscall.Unshare(unix.CLONE_NEWCGROUP); err != nil {
		fatal(err)
	}
	if err := remountFs(); err != nil {
		fatal(err)
	}
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		fatal(err)
	}
	// No need to check args again here. Users can't "accidentally" bypass
	// the parent process and directly reach here. If they do so
	// intentionally, no nice error for them.
	args := os.Args[1:]
	cmd := exec.Cmd{
		Path:   args[0],
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		// Only handle errors that caused the command not to start.
		if _, ok := err.(*exec.ExitError); !ok {
			fatal(err)
		}
	}
	os.Exit(cmd.ProcessState.ExitCode())
}

// setUpChildCgroup checks if the parent has set up cgroup, and continues the
// child-side setup where errors will be fatal.
//
// If successful, we get the following hierarchy:
//   - original
//     subtree_control: on
//     - [UUID]
//       limits: on
//       - worker:
//         procs: [PID]
//     - cgroupName
//       procs: [PPID]
//
// The caller should then unshare its CGROUP namespace to prevent itself or its
// children from modifying limits.
func setUpChildCgroup() {
	c, err := cgroup.New()
	if err != nil || filepath.Base(c.Path) != cgroupName {
		// Cgroup setup was skipped in parent.
		return
	}
	newPath := filepath.Join("../", uuid.NewString())
	c, err = c.MoveToNewSubtree(newPath)
	if err != nil {
		fatal(err)
	}
	_, err = c.MoveToNewSubtree(filepath.Join("worker"))
	if err != nil {
		fatal(err)
	}
	c.EnableLimits()
}

func prepareChildArgs(args []string) ([]string, error) {
	if len(args) == 0 || len(args[0]) == 0 {
		return nil, ErrEmptyArgs
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return nil, err
	}
	return append([]string{specialArg0, path}, args[1:]...), nil
}

func remountFs() error {
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return err
	}
	if err := syscall.Mount("sys", "/sys", "sysfs", 0, ""); err != nil {
		return err
	}
	if err := syscall.Mount("cgroup2", "/sys/fs/cgroup", "cgroup2", 0, ""); err != nil {
		return err
	}
	return nil
}
