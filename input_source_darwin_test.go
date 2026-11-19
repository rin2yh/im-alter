package main

import (
	"fmt"
	"testing"
)

func TestSelectByID(t *testing.T) {
	const (
		id                = "com.apple.keylayout.ABC"
		listRef   uintptr = 1
		queryRef  uintptr = 2
		sourceRef uintptr = 3
	)
	tests := []struct {
		name                  string
		list, query, property uintptr
		status                int32
		wantError             string
		wantSelection         bool
	}{
		{name: "select match", list: listRef, query: queryRef, property: queryRef, status: osStatusOK, wantSelection: true},
		{name: "enumeration failed", wantError: "could not enumerate enabled input sources"},
		{name: "invalid ID", list: listRef, wantError: "invalid UTF-8 input source ID"},
		{name: "missing ID", list: listRef, query: queryRef, wantError: fmt.Sprintf("no enabled input source with ID %q", id)},
		{name: "selection failed", list: listRef, query: queryRef, property: queryRef, status: -50, wantError: fmt.Sprintf("select %q: OSStatus -50", id), wantSelection: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selected := false
			api := inputSourceAPI{
				createSourceList: func(uintptr, bool) uintptr { return tt.list },
				stringFromUTF8:   func(uintptr, string, uint32) uintptr { return tt.query },
				arrayCount:       func(uintptr) int { return 1 },
				arrayValue:       func(uintptr, int) uintptr { return sourceRef },
				sourceProperty:   func(uintptr, uintptr) uintptr { return tt.property },
				equal:            func(left, right uintptr) bool { return left == right },
				release:          func(uintptr) {},
				selectInputSource: func(source uintptr) int32 {
					if source != sourceRef {
						t.Fatalf("selected source = %d, want %d", source, sourceRef)
					}
					selected = true
					return tt.status
				},
			}
			err := api.selectByID(id)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("selectByID: %v", err)
				}
			} else if err == nil || err.Error() != tt.wantError {
				t.Fatalf("selectByID error = %v, want %q", err, tt.wantError)
			}
			if selected != tt.wantSelection {
				t.Errorf("selectInputSource called = %t, want %t", selected, tt.wantSelection)
			}
		})
	}
}
