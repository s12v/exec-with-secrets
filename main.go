package main

import (
	"fmt"
	"github.com/s12v/exec-with-secrets/provider"
	_ "github.com/s12v/exec-with-secrets/provider/awskms"
	_ "github.com/s12v/exec-with-secrets/provider/awssecretsmanager"
	_ "github.com/s12v/exec-with-secrets/provider/awsssm"
	_ "github.com/s12v/exec-with-secrets/provider/azurekeyvault"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("exec-with-secrets:", r)
			os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		fmt.Println("Usage: exec-with-secrets program [args]")
		os.Exit(0)
	}

	path := lookPath(os.Args[1])
	env := provider.Populate(os.Environ())
	_ = syscall.Exec(path, os.Args[1:], env)

	panic("Unable to start " + path)
}

func lookPath(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		panic(err)
	}

	return path
}
