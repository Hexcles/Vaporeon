package worker

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Hexcles/Vaporeon/worker/syncbuffer"
)

var (
	// ErrEmptyArgs is returned by Launch when args or args[0] is empty.
	ErrEmptyArgs = errors.New("worker: empty args")
)

// Init does critical initialization for this library.
//
// It should be called as early as possible: init() is recommended since it
// is guaranteed to run on the startup thread by Go runtime.
//
// Implementation details: in the parent (initial) process, Init() checks and
// sets up cgroups (it does not set any limits on the parent process, but only
// prepares the hierarchy). Later when Launch() starts a new child process, it
// first re-launches itself in a new namespace with a special arg[0], which
// will be detected by Init() in the child process where Init() finishes the
// remaining setup (cgroups, mount points, etc.) and starts the actual job
// process. Init() effectively serves as the PID 1 in the new namespace.
func Init() {
	if os.Args[0] == specialArg0 {
		childInit()
		panic("This should never be reached.")
	} else {
		parentInit()
	}
}

// Job encapsulates a launched process including a private buffer to hold all
// of the job's output.
//
// All exported methods are goroutine-safe. Do not copy this type.
type Job struct {
	// Args record the parameter passed to Launch.
	Args []string
	// Started is guaranteed to be set by Launch.
	Started time.Time

	status JobStatus
	slock  sync.RWMutex

	cmd      *exec.Cmd
	stdout   *syncbuffer.Buffer
	stderr   *syncbuffer.Buffer
	waitOnce sync.Once
}

// JobStatus represents the status of a job. It should be used as values.
type JobStatus struct {
	// Zero value means the job is still running.
	Stopped time.Time
	// If Stopped is zero, this field has no meaning.
	ExitCode int
}

// Launch starts a new job and returns a corresponding *Job instance.
//
// Caller must at least provide a non-empty args[0], which specifies the binary
// to launch (PATH lookup is supported). An error will be returned for empty
// args or invalid/non-existent args[0].
//
// For simplicity, stdin of the launched process is always set to /dev/null.
func Launch(args []string) (*Job, error) {
	childArgs, err := prepareChildArgs(args)
	if err != nil {
		return nil, err
	}
	stdout := syncbuffer.New()
	stderr := syncbuffer.New()
	cmd := &exec.Cmd{
		Path:   "/proc/self/exe",
		Args:   childArgs,
		Stdin:  nil, // /dev/null
		Stdout: stdout,
		Stderr: stderr,
		SysProcAttr: &syscall.SysProcAttr{
			Setpgid: true,
			// CLONE_NEWCGROUP will be unshared from the child init.
			Cloneflags: unix.CLONE_NEWNS | unix.CLONE_NEWPID | unix.CLONE_NEWUSER | unix.CLONE_NEWUTS | unix.CLONE_NEWNET,
			UidMappings: []syscall.SysProcIDMap{
				{HostID: os.Getuid(), ContainerID: 0, Size: 1},
			},
			GidMappings: []syscall.SysProcIDMap{
				{HostID: os.Getgid(), ContainerID: 0, Size: 1},
			},
		},
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	job := &Job{
		Args:    args,
		Started: time.Now(),
		cmd:     cmd,
		stdout:  stdout,
		stderr:  stderr,
	}
	go func() {
		job.waitOnce.Do(job.wait)
	}()
	return job, nil
}

// Stdout returns a Reader to read from the beginning of the job's stdout.
// Call this method multiple times to get multiple independent Readers.
// Read will block when there is no new output and the job is still running;
// EOF is returned only when all output has been read and the job exits.
func (j *Job) Stdout() io.Reader {
	return j.stdout.NewReader()
}

// Stderr is the same as Stdout, but for stderr of the job.
func (j *Job) Stderr() io.Reader {
	return j.stderr.NewReader()
}

// Kill sends SIGKILL to the launched process and waits for it to exit.
//
// Return any error occurred when sending the signal (e.g. process already
// terminated).
func (j *Job) Kill() error {
	if err := j.cmd.Process.Kill(); err != nil {
		return err
	}
	j.waitOnce.Do(j.wait)
	return nil
}

// Status returns a copy of the job's current status.
//
// The status is guaranteed to be available when stdout/stderr returns EOF.
func (j *Job) Status() JobStatus {
	j.slock.RLock()
	defer j.slock.RUnlock()
	return j.status
}

// wait is protected by sync.Once.Do, which drops any error or panic.
func (j *Job) wait() {
	// We will check the exit code later.
	_ = j.cmd.Wait()
	// Set status before closing stdout/err because we guarantee that EOF
	// means the process has stopped and status should be available.
	j.slock.Lock()
	j.status.Stopped = time.Now()
	j.status.ExitCode = j.cmd.ProcessState.ExitCode()
	j.slock.Unlock()
	// Close is crucial; otherwise readers will block for EOF.
	// syncbuffer.Buffer.Close() always succeeds.
	_ = j.stdout.Close()
	_ = j.stderr.Close()
}
