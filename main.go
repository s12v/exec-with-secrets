package main

import (
	"fmt"
	"github.com/s12v/exec-with-secrets/provider"
	_ "github.com/s12v/exec-with-secrets/provider/awskms"
	_ "github.com/s12v/exec-with-secrets/provider/awssecretsmanager"
	_ "github.com/s12v/exec-with-secrets/provider/awsssm"
	_ "github.com/s12v/exec-with-secrets/provider/azurekeyvault"
	"os"
	"syscall"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("exec-with-secrets:", r)
		}
	}()

	if len(os.Args) < 2 {
		fmt.Println("Usage: exec-with-secrets program [args]")
		os.Exit(0)
	}

	env := provider.Populate(os.Environ())
	_ = syscall.Exec(os.Args[1], os.Args[1:], env)

	fmt.Println("exec-with-secrets program: Unable to start", os.Args[1])
	os.Exit(1)
}
