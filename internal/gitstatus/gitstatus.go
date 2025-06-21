// Package gitstatus provides structures and methods to represent and manipulate
// git status information, including branch details and file change statuses.
//
// Structurally, we are looking at something like this:
//
//	StatusInfo
//	├── BranchInfo
//	│   ├── Name
//	│   ├── CommitsAhead
//	│   └── CommitsBehind
//	└── Items
//	    └── StatusItem
//	        ├── ChangeType (enum)
//	        │     ├──> Message():     string
//	        │     ├──> State():       ChangeState (enum)
//	        │     └──> StatusGroup(): StatusGroup (enum)
//	        ├── Path
//	        └── OrigPath
package gitstatus

import (
	"path/filepath"
)

// StatusInfo contains the information about git working tree status that is
// necessary to display the status in a user-friendly way (e.g. git status)
// It includes the branch information and a list of status items.
type StatusInfo struct {
	Branch BranchInfo
	Items  []StatusItem
}

// BranchInfo contains all information needed about the active git branch, as
// well as its status relative to upstream commits.
type BranchInfo struct {
	Name          string // name of the active branch
	CommitsAhead  int    // commit position relative to upstream, e.g. +1
	CommitsBehind int    // commit position relative to upstream, e.g. -3
}

// StatusItem represents a single item of change for a 'git status'.
//
// Note there is not a 1:1 mapping from StatusItem to file paths in underlying
// git status porcelain output; for example, a single path may appear multiple
// times in the output if it has different change types (e.g. staged vs
// unstaged), for example a file that changes that have been staged for commit
// but also has unstaged changes.
type StatusItem struct {
	ChangeType
	Path     string // path relative to the repo root, uses slashes as path separator regardless of OS
	OrigPath string // origin path, e.g. for renamed or copied files, empty otherwise
}

// AbsPath returns the absolute path of the StatusItem based on the git root path
// The path will use the OS native path separator.
//
// See DisplayPath for notes on path normalization and POSIX style.
func (si StatusItem) AbsPath(root string) string {
	root = filepath.ToSlash(filepath.Clean(root))
	return filepath.ToSlash(filepath.Join(root, si.Path))
}

// DisplayPath returns a user-friendly display path for the StatusItem.
// For renamed/copied files, it shows "from -> to" format.
// Paths are shown relative to cwd when possible, otherwise absolute.
// Paths are always in POSIX style for git consistency.
func (si StatusItem) DisplayPath(root, cwd string) string {
	// Normalize all paths to use the same separator style, so that we can
	// reliably calculate relative paths. root and cwd should be in POSIX style
	// for consistency with git repo path in the StatusItem.Path. root probably
	// already is as it comes from git, but cwd comes from os.Getwd() and may be
	// Windows style.  Normalize both anyhow for safety.
	root = filepath.ToSlash(filepath.Clean(root))
	cwd = filepath.ToSlash(filepath.Clean(cwd))

	// relPath converts a repo-relative path to a cwd-relative path
	// or returns the absolute path if relative cannot be calculated.
	//
	// NOTE: both filepath.Rel and filepath.Join will use the OS native separator,
	// so we must convert them back again after those operations.
	relPath := func(repoPath string) string {
		absPath := filepath.ToSlash(filepath.Join(root, repoPath))
		if relPath, err := filepath.Rel(cwd, absPath); err == nil {
			return filepath.ToSlash(relPath)
		}
		return absPath
	}

	// Handle renamed/copied files with "from -> to" format
	if si.OrigPath != "" {
		from := relPath(si.OrigPath)
		to := relPath(si.Path)
		return from + " -> " + to
	}

	return relPath(si.Path)
}

// ChangeType represents the type of change for a path in git status
type ChangeType int

// ChangeType constants represent the different types of changes that can occur in a git status
const (
	// Staged changes
	ChangeStagedModified ChangeType = iota
	ChangeStagedNewFile
	ChangeStagedDeleted
	ChangeStagedRenamed
	ChangeStagedCopied
	ChangeStagedType

	// Unmerged changes (conflicts)
	ChangeUnmergedDeletedBoth
	ChangeUnmergedAddedUs
	ChangeUnmergedDeletedThem
	ChangeUnmergedAddedThem
	ChangeUnmergedDeletedUs
	ChangeUnmergedAddedBoth
	ChangeUnmergedModifiedBoth

	// Unstaged changes
	ChangeUnstagedModified
	ChangeUnstagedDeleted
	ChangeUnstagedType

	// Untracked changes
	ChangeUntracked
)

// changeTypeData maps each changeType to its display information
var changeTypeData = map[ChangeType]changeTypeMetadata{
	ChangeStagedModified:       {msg: "modified", state: ModifiedState, group: Staged},
	ChangeStagedNewFile:        {msg: "new file", state: NewState, group: Staged},
	ChangeStagedDeleted:        {msg: "deleted", state: DeletedState, group: Staged},
	ChangeStagedRenamed:        {msg: "renamed", state: RenamedState, group: Staged},
	ChangeStagedCopied:         {msg: "copied", state: CopiedState, group: Staged},
	ChangeStagedType:           {msg: "typechange", state: TypeChangedState, group: Staged},
	ChangeUnmergedDeletedBoth:  {msg: "both deleted", state: DeletedState, group: Unmerged},
	ChangeUnmergedAddedUs:      {msg: "added by us", state: NewState, group: Unmerged},
	ChangeUnmergedDeletedThem:  {msg: "deleted by them", state: DeletedState, group: Unmerged},
	ChangeUnmergedAddedThem:    {msg: "added by them", state: NewState, group: Unmerged},
	ChangeUnmergedDeletedUs:    {msg: "deleted by us", state: DeletedState, group: Unmerged},
	ChangeUnmergedAddedBoth:    {msg: "both added", state: NewState, group: Unmerged},
	ChangeUnmergedModifiedBoth: {msg: "both modified", state: ModifiedState, group: Unmerged},
	ChangeUnstagedModified:     {msg: "modified", state: ModifiedState, group: Unstaged},
	ChangeUnstagedDeleted:      {msg: "deleted", state: DeletedState, group: Unstaged},
	ChangeUnstagedType:         {msg: "typechange", state: TypeChangedState, group: Unstaged},
	ChangeUntracked:            {msg: "untracked", state: UntrackedState, group: Untracked},
}

// changeTypeMetadata holds the display information for each change type
type changeTypeMetadata struct {
	msg   string
	state ChangeState
	group StatusGroup
}

// ChangeState constants represent the type of a change in a git status output.
// These states describe the nature of the change itself, independent of staging area.
type ChangeState int

const (
	NewState         ChangeState = iota // NewState represents a newly created file
	ModifiedState                       // ModifiedState represents a file with modified content
	DeletedState                        // DeletedState represents a deleted file
	RenamedState                        // RenamedState represents a file that has been renamed
	CopiedState                         // CopiedState represents a file that has been copied
	TypeChangedState                    // TypeChangedState represents a file whose type has changed (e.g., file <-> symlink)
	UntrackedState                      // UntrackedState represents a file not tracked by git
)

// StatusGroup is used to categorize items in git status into groups for
// rendering purposes. Each group represents a different category of the items
// that can appear in the git status output, such as staged, unstaged, merge
// conflicts, and untracked files.
type StatusGroup int

const (
	Staged    StatusGroup = iota // Staged represents changes that are staged for commit
	Unmerged                     // Unmerged represents changes that are in conflict and need resolution
	Unstaged                     // Unstaged represents changes that are not staged for commit
	Untracked                    // Untracked represents files that are not currently tracked by git
)

// Description returns a human-readable description for the StatusGroup,
// typically used for a header or label in the output display.
func (sg StatusGroup) Description() string {
	switch sg {
	case Staged:
		return "Changes to be committed"
	case Unmerged:
		return "Unmerged paths"
	case Unstaged:
		return "Changes not staged for commit"
	case Untracked:
		return "Untracked files"
	default:
		panic("invalid status group")
	}
}

// Message returns the display message for the change type
func (ct ChangeType) Message() string {
	if info, ok := changeTypeData[ct]; ok {
		return info.msg
	}
	panic("invalid change type")
}

// State returns the change state associated with the change type
func (ct ChangeType) State() ChangeState {
	if info, ok := changeTypeData[ct]; ok {
		return info.state
	}
	panic("invalid change type")
}

// StatusGroup returns the status group for the change type
func (ct ChangeType) StatusGroup() StatusGroup {
	if info, ok := changeTypeData[ct]; ok {
		return info.group
	}
	panic("invalid change type")
}
