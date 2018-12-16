package main

import (
	"fmt"
	"github.com/s12v/secure-exec/provider"
	_ "github.com/s12v/secure-exec/provider/awskms"
	_ "github.com/s12v/secure-exec/provider/awssecretsmanager"
	_ "github.com/s12v/secure-exec/provider/awssecretsmanager"
	_ "github.com/s12v/secure-exec/provider/awsssm"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: secure-exec program [args]")
		os.Exit(0)
	}

	env := provider.Populate(os.Environ())
	syscall.Exec(os.Args[1], os.Args[1:], env);

	fmt.Printf("Unable to start %v", os.Args[1])
	os.Exit(1)
}
