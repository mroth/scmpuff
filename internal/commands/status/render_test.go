package status

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	updateGolden  = flag.Bool("update", false, "update the golden files of this test")
	clobberGolden = flag.Bool("clobber", false, "clobber all golden files in testdata and exit")
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *clobberGolden {
		goldenClobberFiles()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestRenderer_Display(t *testing.T) {
	testCases := []struct {
		name string
		info StatusInfo
	}{
		{
			name: "empty",
			info: StatusInfo{
				BranchInfo{Name: "main", CommitsAhead: 0, CommitsBehind: 0},
				nil,
			},
		},
		{
			name: "with_branch_ahead",
			info: StatusInfo{
				BranchInfo{Name: "feature", CommitsAhead: 3, CommitsBehind: 0},
				nil,
			},
		},
		{
			name: "with_staged_files",
			info: StatusInfo{
				BranchInfo{Name: "main", CommitsAhead: 0, CommitsBehind: 0},
				[]StatusItem{
					{ChangeType: ChangeStagedNewFile, FileAbsPath: "/path/to/new.go", FileRelPath: "new.go"},
					{ChangeType: ChangeStagedNewFile, FileAbsPath: "/path/to/new_b.go", FileRelPath: "new_b.go"},
					{ChangeType: ChangeStagedModified, FileAbsPath: "/path/to/changed.go", FileRelPath: "changed.go"}},
			},
		},
		{
			name: "complex_mix",
			info: StatusInfo{
				BranchInfo{Name: "feature", CommitsAhead: 2, CommitsBehind: 1},
				[]StatusItem{
					{ChangeType: ChangeStagedNewFile, FileAbsPath: "/path/to/new.go", FileRelPath: "new.go"},
					{ChangeType: ChangeStagedNewFile, FileAbsPath: "/path/to/new_b.go", FileRelPath: "new_b.go"},
					{ChangeType: ChangeUnstagedModified, FileAbsPath: "/path/to/modified.go", FileRelPath: "modified.go"},
					{ChangeType: ChangeUntracked, FileAbsPath: "/path/to/untracked.go", FileRelPath: "untracked.go"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var optionCases = []struct {
				name                string
				includeParseData    bool
				includeStatusOutput bool
			}{
				{name: "parsedata.txt", includeParseData: true, includeStatusOutput: false},
				{name: "display.ansi", includeParseData: false, includeStatusOutput: true},
			}

			for _, oc := range optionCases {
				t.Run(oc.name, func(t *testing.T) {
					renderer, err := NewRenderer(&tc.info)
					if err != nil {
						t.Fatalf("NewRenderer() error: %v", err)
					}

					var buf bytes.Buffer
					err = renderer.Display(&buf, oc.includeParseData, oc.includeStatusOutput)
					if err != nil {
						t.Fatalf("Display() error: %v", err)
					}

					goldenFile := fmt.Sprintf("statuslist-%s.%s", tc.name, oc.name)
					goldenCompareFile(t, goldenFile, buf.Bytes(), *updateGolden)
				})
			}
		})
	}
}

func goldenCompareFile(t *testing.T, filename string, actual []byte, update bool) {
	t.Helper()

	goldenPath := "testdata/" + filename + ".golden"

	if update {
		err := os.MkdirAll(filepath.Dir(goldenPath), 0755)
		if err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}

		err = os.WriteFile(goldenPath, actual, 0644)
		if err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
		t.Logf("updated golden file: %s [%v bytes]", goldenPath, len(actual))
	}

	goldenData, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v\nRun with -update to create it", goldenPath, err)
	}

	if !bytes.Equal(goldenData, actual) {
		t.Errorf("actual doesn't match golden file %s\nExpected:\n%s\nActual:\n%s",
			goldenPath, goldenData, actual)
	}
}

func goldenClobberFiles() {
	files, err := filepath.Glob("testdata/*.golden")
	if err != nil {
		panic(err) // only possible if the glob pattern is invalid
	}
	log.Printf("ℹ️ found %d golden files in testdata", len(files))

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			log.Fatalf("❌ failed to clobber golden file %s: %v", file, err)
		}
		log.Printf("♻️ removed golden file %s", file)
	}
}
