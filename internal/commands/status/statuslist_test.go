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

func TestStatusList_Display(t *testing.T) {
	testCases := []struct {
		name       string
		statusList *StatusList
	}{
		{
			name: "empty",
			statusList: createTestStatusList(
				&BranchInfo{name: "main", ahead: 0, behind: 0},
				nil,
			),
		},
		{
			name: "with_branch_ahead",
			statusList: createTestStatusList(
				&BranchInfo{name: "feature", ahead: 3, behind: 0},
				nil,
			),
		},
		{
			name: "with_staged_files",
			statusList: createTestStatusList(
				&BranchInfo{name: "main", ahead: 0, behind: 0},
				[]StatusItem{
					{msg: "  new file", col: neu, group: Staged, fileAbsPath: "/path/to/new.go", fileRelPath: "new.go"},
					{msg: "  new file", col: neu, group: Staged, fileAbsPath: "/path/to/new_b.go", fileRelPath: "new_b.go"},
					{msg: "  modified", col: mod, group: Staged, fileAbsPath: "/path/to/changed.go", fileRelPath: "changed.go"},
				},
			),
		},
		{
			name: "complex_mix",
			statusList: createTestStatusList(
				&BranchInfo{name: "feature", ahead: 2, behind: 1},
				[]StatusItem{
					{msg: "  new file", col: neu, group: Staged, fileAbsPath: "/path/to/new.go", fileRelPath: "new.go"},
					{msg: "  new file", col: neu, group: Staged, fileAbsPath: "/path/to/new_b.go", fileRelPath: "new_b.go"},
					{msg: "  modified", col: mod, group: Unstaged, fileAbsPath: "/path/to/modified.go", fileRelPath: "modified.go"},
					{msg: " untracked", col: unt, group: Untracked, fileAbsPath: "/path/to/untracked.go", fileRelPath: "untracked.go"},
				},
			),
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
					var buf bytes.Buffer
					err := tc.statusList.Display(&buf, oc.includeParseData, oc.includeStatusOutput)
					if err != nil {
						t.Fatalf("error calling Display: %v", err)
					}

					goldenFile := fmt.Sprintf("statuslist-%s.%s", tc.name, oc.name)
					goldenCompareFile(t, goldenFile, buf.Bytes(), *updateGolden)
				})
			}
		})
	}
}

// Helper function to create a test StatusList
func createTestStatusList(branch *BranchInfo, items []StatusItem) *StatusList {
	sl := NewStatusList()
	sl.branch = branch

	for _, si := range items {
		sl.Add(&si)
	}

	return sl
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
