package status

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/kmatt/scmpuff/internal/gitstatus"
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
		name      string
		info      gitstatus.StatusInfo
		root, cwd string
	}{
		{
			name: "empty",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "main", CommitsAhead: 0, CommitsBehind: 0},
				Items:  nil,
			},
		},
		{
			name: "with_branch_ahead",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "feature", CommitsAhead: 3, CommitsBehind: 0},
				Items:  nil,
			},
		},
		{
			name: "with_staged_files",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "main", CommitsAhead: 0, CommitsBehind: 0},
				Items: []gitstatus.StatusItem{
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new.go"},
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new_b.go"},
					{ChangeType: gitstatus.ChangeStagedModified, Path: "changed.go"}},
			},
			root: "/path/to",
			cwd:  "/path/to",
		},
		{
			name: "complex_mix",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "feature", CommitsAhead: 2, CommitsBehind: 1},
				Items: []gitstatus.StatusItem{
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new.go"},
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new_b.go"},
					{ChangeType: gitstatus.ChangeUnstagedModified, Path: "modified.go"},
					{ChangeType: gitstatus.ChangeUntracked, Path: "untracked.go"},
				},
			},
			root: "/path/to",
			cwd:  "/path/to",
		},
		{
			// longer list of changes (more than 10), unicode, some emoji, copy, rename, delete
			name: "longlist",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "techdebt", CommitsAhead: 42, CommitsBehind: 1123},
				Items: []gitstatus.StatusItem{
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new_a.php"},
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new_b.php"},
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new_c.php"},
					{ChangeType: gitstatus.ChangeStagedNewFile, Path: "new_d.php"},
					{ChangeType: gitstatus.ChangeUnstagedModified, Path: "modified1.php"},
					{ChangeType: gitstatus.ChangeUnstagedModified, Path: "modified2.php"},
					{ChangeType: gitstatus.ChangeUnstagedModified, Path: "修改后的文件.php"},
					{ChangeType: gitstatus.ChangeUntracked, Path: "untracked file with spaces.txt"},
					{ChangeType: gitstatus.ChangeStagedRenamed, Path: "tests/disabled", OrigPath: "tests/flakey"},
					{ChangeType: gitstatus.ChangeStagedRenamed, Path: "docs/SECURITY.md", OrigPath: "SECURITY.md"},
					{ChangeType: gitstatus.ChangeStagedCopied, Path: "metoo", OrigPath: "me"},
					{ChangeType: gitstatus.ChangeUnstagedDeleted, Path: "👻.go"},
				},
			},
			root: "/Users/bobbytables/code",
			cwd:  "/Users/bobbytables/code",
		},
		{
			name: "subdirectory",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "feature", CommitsAhead: 0, CommitsBehind: 13},
				Items: []gitstatus.StatusItem{
					{ChangeType: gitstatus.ChangeStagedRenamed, Path: "projects/snw", OrigPath: "projects/ds9"},
					{ChangeType: gitstatus.ChangeStagedRenamed, Path: "projects/warpcore/CONFIDENTIAL.md", OrigPath: "projects/warpcore/SporeDriveSchematics.md"},
					{ChangeType: gitstatus.ChangeStagedDeleted, Path: "docs/wolf 359 was an inside job.txt"},
				},
			},
			root: "/home/starfleet/src",
			cwd:  "/home/starfleet/src/projects/warpcore",
		},
		{
			name: "unmerged_conflicts",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "merge-conflict", CommitsAhead: 0, CommitsBehind: 0},
				Items: []gitstatus.StatusItem{
					{ChangeType: gitstatus.ChangeUnmergedDeletedBoth, Path: "deleted_by_both.txt"},
					{ChangeType: gitstatus.ChangeUnmergedAddedUs, Path: "added_by_us.txt"},
					{ChangeType: gitstatus.ChangeUnmergedDeletedThem, Path: "deleted_by_them.txt"},
					{ChangeType: gitstatus.ChangeUnmergedAddedThem, Path: "added_by_them.txt"},
					{ChangeType: gitstatus.ChangeUnmergedDeletedUs, Path: "deleted_by_us.txt"},
					{ChangeType: gitstatus.ChangeUnmergedAddedBoth, Path: "added_by_both.txt"},
					{ChangeType: gitstatus.ChangeUnmergedModifiedBoth, Path: "modified_by_both.txt"},
				},
			},
			root: "/path/to/repo",
			cwd:  "/path/to/repo",
		},
		{
			name: "type_changes",
			info: gitstatus.StatusInfo{
				Branch: gitstatus.BranchInfo{Name: "type-change", CommitsAhead: 1, CommitsBehind: 0},
				Items: []gitstatus.StatusItem{
					{ChangeType: gitstatus.ChangeStagedType, Path: "staged_typechange.txt"},
					{ChangeType: gitstatus.ChangeUnstagedType, Path: "unstaged_typechange.txt"},
				},
			},
			root: "/path/to/repo",
			cwd:  "/path/to/repo",
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
					renderer, err := NewRenderer(&tc.info, tc.root, tc.cwd)
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
