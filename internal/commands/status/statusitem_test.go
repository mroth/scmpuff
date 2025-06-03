package status

import (
	"runtime"
	"testing"
)

func TestStatusItem(t *testing.T) {
	testcases := []struct {
		name            string
		item            StatusItem
		root            string
		wd              string
		wantAbsPath     string
		wantDispPath    string
		windowsOnlyTest bool // flag to indicate if this test is only relevant for Windows
	}{
		{
			name:         "everything in root",
			item:         StatusItem{Path: "a.txt"},
			root:         "/tmp/foo",
			wd:           "/tmp/foo",
			wantAbsPath:  "/tmp/foo/a.txt",
			wantDispPath: "a.txt",
		},
		{
			name:         "change in subdir",
			item:         StatusItem{Path: "bar/c.txt"},
			root:         "/tmp/foo",
			wd:           "/tmp/foo",
			wantAbsPath:  "/tmp/foo/bar/c.txt",
			wantDispPath: "bar/c.txt",
		},
		{
			name:         "change in parent to wd",
			item:         StatusItem{Path: "a.txt"},
			root:         "/tmp/foo",
			wd:           "/tmp/foo/bar",
			wantAbsPath:  "/tmp/foo/a.txt",
			wantDispPath: "../a.txt",
		},
		{
			name:         "handle trailing slashes",
			item:         StatusItem{Path: "a.txt"},
			root:         "/tmp/foo/",
			wd:           "/tmp/foo/bar/",
			wantAbsPath:  "/tmp/foo/a.txt",
			wantDispPath: "../a.txt",
		},
		{
			name:         "copyrename paths",
			item:         StatusItem{Path: "b.txt", OrigPath: "a.txt"},
			root:         "/tmp/foo",
			wd:           "/tmp/foo",
			wantAbsPath:  "/tmp/foo/b.txt",
			wantDispPath: "a.txt -> b.txt",
		},
		{
			name:         "copyrename paths in subdir",
			item:         StatusItem{Path: "bar/b.txt", OrigPath: "bar/a.txt"},
			root:         "/tmp/foo",
			wd:           "/tmp/foo",
			wantAbsPath:  "/tmp/foo/bar/b.txt",
			wantDispPath: "bar/a.txt -> bar/b.txt",
		},
		{
			name:         "copyrename paths in parent dir",
			item:         StatusItem{Path: "b.txt", OrigPath: "a.txt"},
			root:         "/tmp/foo",
			wd:           "/tmp/foo/bar",
			wantAbsPath:  "/tmp/foo/b.txt",
			wantDispPath: "../a.txt -> ../b.txt",
		},
		{
			name:            "windows style paths",
			item:            StatusItem{Path: "bar/a.txt"}, // StatusItem Paths are always POSIX style, like git
			root:            `C:/Users/Bob/foo`,            // root will likely be POSIX style, since being retrieved from git
			wd:              `C:\Users\Bob\foo`,            // wd is Windows style, e.g. from os.Getwd()
			wantAbsPath:     `C:/Users/Bob/foo/bar/a.txt`,  // AbsPath should be POSIX style (?)
			wantDispPath:    `bar/a.txt`,                   // DisplayPath should be POSIX style for git consistency
			windowsOnlyTest: true,                          // this test is only relevant for Windows
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.windowsOnlyTest && runtime.GOOS != "windows" {
				t.Skip("Skipping Windows-only test on non-Windows platform")
			}

			gotAbsPath := tc.item.AbsPath(tc.root)
			if gotAbsPath != tc.wantAbsPath {
				t.Errorf("StatusItem.AbsPath() = %v, want %v", gotAbsPath, tc.wantAbsPath)
			}

			gotRelPath := tc.item.DisplayPath(tc.root, tc.wd)
			if gotRelPath != tc.wantDispPath {
				t.Errorf("StatusItem.DisplayPath() = %v, want %v", gotRelPath, tc.wantDispPath)
			}
		})
	}
}
