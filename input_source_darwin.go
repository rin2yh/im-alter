package main

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

func selectSource(id string) error {
	api, err := loadInputSourceAPI()
	if err != nil {
		return err
	}

	sources := api.createSourceList(0, false)
	if sources == 0 {
		return fmt.Errorf("could not enumerate enabled input sources")
	}
	defer api.release(sources)

	source, err := api.findByID(sources, id)
	if err != nil {
		return err
	}
	if status := api.selectInputSource(source); status != 0 {
		return fmt.Errorf("select %q: OSStatus %d", id, status)
	}
	return nil
}

type inputSourceAPI struct {
	idKey             uintptr
	createSourceList  func(uintptr, bool) uintptr
	sourceProperty    func(uintptr, uintptr) uintptr
	selectInputSource func(uintptr) int32
	stringFromUTF8    func(uintptr, string, uint32) uintptr
	arrayCount        func(uintptr) int
	arrayValue        func(uintptr, int) uintptr
	equal             func(uintptr, uintptr) bool
	release           func(uintptr)
}

func (api *inputSourceAPI) findByID(sources uintptr, id string) (uintptr, error) {
	const utf8Encoding = 0x08000100
	want := api.stringFromUTF8(0, id, utf8Encoding)
	if want == 0 {
		return 0, fmt.Errorf("invalid UTF-8 input source ID")
	}
	defer api.release(want)

	for i := range api.arrayCount(sources) {
		source := api.arrayValue(sources, i)
		if sourceID := api.sourceProperty(source, api.idKey); sourceID != 0 && api.equal(sourceID, want) {
			return source, nil
		}
	}
	return 0, fmt.Errorf("no enabled input source with ID %q", id)
}

func loadInputSourceAPI() (*inputSourceAPI, error) {
	const frameworks = "/System/Library/Frameworks/"
	carbon, err := purego.Dlopen(frameworks+"Carbon.framework/Carbon", purego.RTLD_NOW)
	if err != nil {
		return nil, err
	}
	cf, err := purego.Dlopen(frameworks+"CoreFoundation.framework/CoreFoundation", purego.RTLD_NOW)
	if err != nil {
		return nil, err
	}
	key, err := purego.Dlsym(carbon, "kTISPropertyInputSourceID")
	if err != nil {
		return nil, err
	}

	// The exported variable contains the CFStringRef used as the property key.
	api := &inputSourceAPI{idKey: *(*uintptr)(unsafe.Pointer(key))}
	purego.RegisterLibFunc(&api.createSourceList, carbon, "TISCreateInputSourceList")
	purego.RegisterLibFunc(&api.sourceProperty, carbon, "TISGetInputSourceProperty")
	purego.RegisterLibFunc(&api.selectInputSource, carbon, "TISSelectInputSource")
	purego.RegisterLibFunc(&api.stringFromUTF8, cf, "CFStringCreateWithCString")
	purego.RegisterLibFunc(&api.arrayCount, cf, "CFArrayGetCount")
	purego.RegisterLibFunc(&api.arrayValue, cf, "CFArrayGetValueAtIndex")
	purego.RegisterLibFunc(&api.equal, cf, "CFEqual")
	purego.RegisterLibFunc(&api.release, cf, "CFRelease")
	return api, nil
}
