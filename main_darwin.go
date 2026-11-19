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
		fmt.Fprintln(os.Stderr, "usage: im-alter {abc|INPUT_SOURCE_ID}")
		os.Exit(2)
	}
	id := os.Args[1]
	if id == "abc" {
		id = "com.apple.keylayout.ABC"
	}
	if err := selectSource(id); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
