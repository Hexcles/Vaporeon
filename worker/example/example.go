package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Hexcles/Vaporeon/worker"
)

func init() {
	worker.Init()
}

func run(args ...string) {
	fmt.Printf("* Running: %v\n", args)
	job, err := worker.Launch(args)
	if err != nil {
		panic(err)
	}
	stdout, err := ioutil.ReadAll(job.Stdout())
	stderr, err := ioutil.ReadAll(job.Stderr())
	if err != nil {
		panic(err)
	}
	fmt.Printf("* Exit code: %d\n", *job.ExitCode)
	fmt.Println("* Stdout:")
	fmt.Print(string(stdout))
	fmt.Println("* Stderr:")
	fmt.Print(string(stderr))
}

func main() {
	if len(os.Args) > 1 {
		run(os.Args[1:]...)
		return
	}
	run("whoami")
	run("hostname")
	run("ls", "-l", "/sys/fs/cgroup")
	run("ip", "link")
	run("ps", "-ef")
}