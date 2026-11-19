// Command im-alter selects an enabled macOS input source by its identifier.
package main

import (
	"fmt"
	"os"
	"runtime"
)

func init() {
	// Text Input Source Services needs the process's main thread.
	runtime.LockOSThread()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: im-alter INPUT_SOURCE_ID")
		os.Exit(2)
	}
	if err := selectSource(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
