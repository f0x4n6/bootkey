// Extract the BootKey from an offline system hive.
//
// Usage:
//
//	bootkey system
//
// The arguments are:
//
//	system
//	    The system registry hive (required).
package main

import (
	"fmt"
	"os"

	"go.foxforensics.dev/bootkey/pkg/bootkey"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" {
		_, _ = fmt.Fprintln(os.Stderr, "usage: bootkey system")
		os.Exit(2)
	}

	key, err := bootkey.ReadFile(os.Args[1])

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%x\n", key)
}
