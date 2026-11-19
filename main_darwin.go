// Command im-alter selects an enabled macOS input source by its identifier.
package main

import (
	"fmt"
	"os"
	"runtime"
)

var inputSourceAliases = map[string]string{
	"abc":      "com.apple.keylayout.ABC",
	"hiragana": "com.apple.inputmethod.Kotoeri.RomajiTyping.Japanese",
	"us":       "com.apple.keylayout.US",
}

func init() {
	// Text Input Source Services needs the process's main thread.
	runtime.LockOSThread()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: im-alter {ALIAS|INPUT_SOURCE_ID}")
		os.Exit(2)
	}
	id := os.Args[1]
	if sourceID, ok := inputSourceAliases[id]; ok {
		id = sourceID
	}
	if err := selectSource(id); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
