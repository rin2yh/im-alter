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
	type input struct {
		id           string
		sourceList   uintptr
		queryString  uintptr
		sourceID     uintptr
		selectStatus int32
	}
	type want struct {
		err      string
		selected bool
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name:  "select match",
			input: input{id: id, sourceList: listRef, queryString: queryRef, sourceID: queryRef, selectStatus: osStatusOK},
			want:  want{selected: true},
		},
		{
			name:  "enumeration failed",
			input: input{id: id},
			want:  want{err: "could not enumerate enabled input sources"},
		},
		{
			name:  "invalid ID",
			input: input{id: "\xff", sourceList: listRef},
			want:  want{err: "invalid UTF-8 input source ID"},
		},
		{
			name:  "missing ID",
			input: input{id: id, sourceList: listRef, queryString: queryRef},
			want:  want{err: fmt.Sprintf("no enabled input source with ID %q", id)},
		},
		{
			name:  "selection failed",
			input: input{id: id, sourceList: listRef, queryString: queryRef, sourceID: queryRef, selectStatus: -50},
			want:  want{err: fmt.Sprintf("select %q: OSStatus -50", id), selected: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selected := false
			api := inputSourceAPI{
				createSourceList: func(uintptr, bool) uintptr { return tt.input.sourceList },
				stringFromUTF8:   func(uintptr, string, uint32) uintptr { return tt.input.queryString },
				arrayCount:       func(uintptr) int { return 1 },
				arrayValue:       func(uintptr, int) uintptr { return sourceRef },
				sourceProperty:   func(uintptr, uintptr) uintptr { return tt.input.sourceID },
				equal:            func(left, right uintptr) bool { return left == right },
				release:          func(uintptr) {},
				selectInputSource: func(source uintptr) int32 {
					if source != sourceRef {
						t.Fatalf("selected source = %d, want %d", source, sourceRef)
					}
					selected = true
					return tt.input.selectStatus
				},
			}
			err := api.selectByID(tt.input.id)
			if tt.want.err == "" {
				if err != nil {
					t.Fatalf("selectByID: %v", err)
				}
			} else if err == nil || err.Error() != tt.want.err {
				t.Fatalf("selectByID error = %v, want %q", err, tt.want.err)
			}
			if selected != tt.want.selected {
				t.Errorf("selectInputSource called = %t, want %t", selected, tt.want.selected)
			}
		})
	}
}
