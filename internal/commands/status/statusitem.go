package status

import (
	"path/filepath"
)

// StatusItem represents a single processed item of change from a 'git status'
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
	ChangeUnmergedDeletedBoth ChangeType = iota
	ChangeUnmergedAddedUs
	ChangeUnmergedDeletedThem
	ChangeUnmergedAddedThem
	ChangeUnmergedDeletedUs
	ChangeUnmergedAddedBoth
	ChangeUnmergedModifiedBoth

	ChangeUntracked

	ChangeStagedModified
	ChangeStagedNewFile
	ChangeStagedDeleted
	ChangeStagedRenamed
	ChangeStagedCopied
	ChangeStagedType

	ChangeUnstagedModified
	ChangeUnstagedDeleted
	ChangeUnstagedType
)

// changeTypeData maps each changeType to its display information
var changeTypeData = map[ChangeType]changeTypeInfo{
	ChangeUnmergedDeletedBoth:  {msg: "   both deleted", state: DeletedState, group: Unmerged},
	ChangeUnmergedAddedUs:      {msg: "    added by us", state: NewState, group: Unmerged},
	ChangeUnmergedDeletedThem:  {msg: "deleted by them", state: DeletedState, group: Unmerged},
	ChangeUnmergedAddedThem:    {msg: "  added by them", state: NewState, group: Unmerged},
	ChangeUnmergedDeletedUs:    {msg: "  deleted by us", state: DeletedState, group: Unmerged},
	ChangeUnmergedAddedBoth:    {msg: "     both added", state: NewState, group: Unmerged},
	ChangeUnmergedModifiedBoth: {msg: "  both modified", state: ModifiedState, group: Unmerged},
	ChangeUntracked:            {msg: " untracked", state: UntrackedState, group: Untracked},
	ChangeStagedModified:       {msg: "  modified", state: ModifiedState, group: Staged},
	ChangeStagedNewFile:        {msg: "  new file", state: NewState, group: Staged},
	ChangeStagedDeleted:        {msg: "   deleted", state: DeletedState, group: Staged},
	ChangeStagedRenamed:        {msg: "   renamed", state: RenamedState, group: Staged},
	ChangeStagedCopied:         {msg: "    copied", state: CopiedState, group: Staged},
	ChangeStagedType:           {msg: "typechange", state: TypeChangedState, group: Staged},
	ChangeUnstagedModified:     {msg: "  modified", state: ModifiedState, group: Unstaged},
	ChangeUnstagedDeleted:      {msg: "   deleted", state: DeletedState, group: Unstaged},
	ChangeUnstagedType:         {msg: "typechange", state: TypeChangedState, group: Unstaged},
}

// changeTypeInfo holds the display information for each change type
type changeTypeInfo struct {
	msg   string
	state changeState
	group StatusGroup
}

type changeState int

const (
	NewState changeState = iota
	ModifiedState
	DeletedState
	UntrackedState
	RenamedState
	CopiedState
	TypeChangedState
)

// StatusGroup encapsulates constants for mapping group status
type StatusGroup int

// constants representing an enum of all possible StatusGroups
const (
	Staged StatusGroup = iota
	Unmerged
	Unstaged
	Untracked
)

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
		return "Unknown group"
	}
}

// Message returns the display message for the change type
func (ct ChangeType) Message() string {
	if info, ok := changeTypeData[ct]; ok {
		return info.msg
	}
	panic("invalid change type")
}

// State returns the change state for the change type
func (ct ChangeType) state() changeState {
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
