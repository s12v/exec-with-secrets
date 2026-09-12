package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	_ "github.com/s12v/exec-with-secrets/provider/awskms"
	_ "github.com/s12v/exec-with-secrets/provider/awssecretsmanager"
	_ "github.com/s12v/exec-with-secrets/provider/awsssm"

	"github.com/s12v/exec-with-secrets/provider"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "exec-with-secrets:", r)
			os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: exec-with-secrets program [args]")
		os.Exit(2)
	}

	path := lookPath(os.Args[1])
	env := provider.Populate(os.Environ())
	if err := syscall.Exec(path, os.Args[1:], env); err != nil {
		panic("Unable to start " + path + ": " + err.Error())
	}
}

func lookPath(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		panic(err)
	}

	return path
}
