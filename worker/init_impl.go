package worker

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

const specialArg0 = "<init>"
const hostname = "vaporeon"

func childInit() {
	runtime.LockOSThread()
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

func fatal(err error) {
	log.Fatalf("Init error: %s", err)
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
