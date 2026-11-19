// Command im-alter selects an enabled macOS input source by its identifier.
package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
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

func selectSource(id string) error {
	const frameworks = "/System/Library/Frameworks/"
	carbon, err := purego.Dlopen(frameworks+"Carbon.framework/Carbon", purego.RTLD_NOW)
	if err != nil {
		return err
	}
	cf, err := purego.Dlopen(frameworks+"CoreFoundation.framework/CoreFoundation", purego.RTLD_NOW)
	if err != nil {
		return err
	}

	var list func(uintptr, bool) uintptr
	var property func(uintptr, uintptr) uintptr
	var selectInput func(uintptr) int32
	purego.RegisterLibFunc(&list, carbon, "TISCreateInputSourceList")
	purego.RegisterLibFunc(&property, carbon, "TISGetInputSourceProperty")
	purego.RegisterLibFunc(&selectInput, carbon, "TISSelectInputSource")

	var stringFromUTF8 func(uintptr, string, uint32) uintptr
	var count func(uintptr) int
	var item func(uintptr, int) uintptr
	var equal func(uintptr, uintptr) bool
	var release func(uintptr)
	purego.RegisterLibFunc(&stringFromUTF8, cf, "CFStringCreateWithCString")
	purego.RegisterLibFunc(&count, cf, "CFArrayGetCount")
	purego.RegisterLibFunc(&item, cf, "CFArrayGetValueAtIndex")
	purego.RegisterLibFunc(&equal, cf, "CFEqual")
	purego.RegisterLibFunc(&release, cf, "CFRelease")

	key, err := purego.Dlsym(carbon, "kTISPropertyInputSourceID")
	if err != nil {
		return err
	}
	// The exported variable contains the CFStringRef used as the property key.
	idKey := *(*uintptr)(unsafe.Pointer(key))
	const utf8Encoding = 0x08000100
	want := stringFromUTF8(0, id, utf8Encoding)
	if want == 0 {
		return fmt.Errorf("invalid UTF-8 input source ID")
	}
	defer release(want)
	sources := list(0, false)
	if sources == 0 {
		return fmt.Errorf("could not enumerate enabled input sources")
	}
	defer release(sources)
	for i := range count(sources) {
		source := item(sources, i)
		if got := property(source, idKey); got != 0 && equal(got, want) {
			if status := selectInput(source); status != 0 {
				return fmt.Errorf("select %q: OSStatus %d", id, status)
			}
			return nil
		}
	}
	return fmt.Errorf("no enabled input source with ID %q", id)
}
