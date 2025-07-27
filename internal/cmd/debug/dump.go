package debug

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const (
	tmpDirPattern   = "scmpuff-dumptool-"
	timestampFormat = "20060102150405"
)

type dumpOptions struct {
	wantSummary   bool
	wantPorcelain bool
	wantArchive   bool
	verbose       bool
}

var opts dumpOptions

// NewDumpCmd creates and returns the dump command
func NewDumpCmd() *cobra.Command {
	dumpCmd := &cobra.Command{
		Use:   "dump",
		Short: "Dump git status information for debugging",
		Long: `Dump git status information and porcelain files in a machine-readable format for debugging or analysis purposes.

This tool captures various git status formats and metadata to help debug issues with scmpuff.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true // silence usage-on-error after args processed
			return run(&opts)
		},
		Example: `
# Basic usage (summary only)
scmpuff debug dump

# Preserves porcelain files locally for inspection
scmpuff debug dump --porcelain

# Create an archive of all data
scmpuff debug dump --archive

# Disable summary
scmpuff debug dump --summary=false

# Combined flags
scmpuff debug dump --porcelain --summary=false -v`,
	}

	dumpCmd.Flags().BoolVar(&opts.wantSummary, "summary", true, "summary of status for debugging")
	dumpCmd.Flags().BoolVar(&opts.wantPorcelain, "porcelain", false, "keep porcelain files for test or debug purposes")
	dumpCmd.Flags().BoolVar(&opts.wantArchive, "archive", false, "create compressed archive of data")
	dumpCmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "enable verbose output")

	return dumpCmd
}

func run(opts *dumpOptions) error {
	// Step 1: Gather environment information
	report, err := collectReport()
	if err != nil {
		return err
	}

	// Step 2: Display summary if requested
	if opts.wantSummary {
		fmt.Println(report.Summary())
	}

	// Step 3: Handle file operations (porcelain files and archive)
	if opts.wantPorcelain || opts.wantArchive {
		// Create a temporary directory for porcelain files.
		tmpdir, err := os.MkdirTemp("", tmpDirPattern)
		if err != nil {
			return fmt.Errorf("could not make tmpdir: %w", err)
		}

		// If user specifically asked for porcelain files, keep them around afterwards,
		// otherwise clean up tmpdir after we are done.
		if !opts.wantPorcelain {
			defer func() {
				err := os.RemoveAll(tmpdir)
				if opts.verbose {
					log.Printf("VERBOSE: removed tmpdir %s, error = %v", tmpdir, err)
				}
			}()
		}

		// Write a summary file if summary was requested
		if opts.wantSummary {
			err := writeSummaryFile(tmpdir, report.Slug(), report.Summary(), opts.verbose)
			if err != nil {
				return fmt.Errorf("failed to write summary file: %w", err)
			}
		}

		// Write porcelain files
		err = createPorcelainSamples(tmpdir, report.Slug(), opts.verbose)
		if err != nil {
			return fmt.Errorf("failed to write porcelain samples: %w", err)
		}
		if opts.wantPorcelain { // Notify user of porcelain file location if they specifically requested them
			log.Printf("Porcelain files saved to: %s", tmpdir)
		}

		// Create archive if requested
		if opts.wantArchive {
			if err := createArchive(tmpdir, report.Slug()); err != nil {
				return fmt.Errorf("failed to create archive: %w", err)
			}
		}
	}

	return nil
}

type EnvironmentReport struct {
	cwd        string // OS working directory
	user       string // OS username
	gitVersion string // git version string
	gitRoot    string // git project root, determined by 'git rev-parse --show-toplevel'
	cdupPath   string // path from 'git rev-parse --show-cdup'
	statusData []byte
	timestamp  time.Time
}

func collectReport() (*EnvironmentReport, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve current working directory: %w", err)
	}

	gitVersionBytes, err := exec.Command("git", "version").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to determine git version: %w", err)
	}

	toplevelBytes, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to determine git project root: %w", err)
	}

	cdupBytes, err := exec.Command("git", "rev-parse", "--show-cdup").Output()
	if err != nil {
		return nil, fmt.Errorf("failed on call to --show-cdup: %w", err)
	}

	statusBytes, err := exec.Command("git", "status", "--branch", "--porcelain=v1").Output()
	if err != nil {
		return nil, fmt.Errorf("git status error: %w", err)
	}

	userFromOS, err := user.Current()
	user := "anonymous"
	if err == nil {
		user = userFromOS.Username
	}

	return &EnvironmentReport{
		cwd:        wd,
		user:       user,
		gitVersion: string(bytes.TrimSpace(gitVersionBytes)),
		gitRoot:    string(bytes.TrimSpace(toplevelBytes)),
		cdupPath:   string(bytes.TrimSpace(cdupBytes)),
		statusData: statusBytes,
		timestamp:  time.Now(),
	}, nil
}

// Slug returns a unique ID slug based on this environment report, for use in filenames
func (r *EnvironmentReport) Slug() string {
	return fmt.Sprintf("%s-%s-%s",
		r.user,
		filepath.Base(r.gitRoot),
		r.timestamp.Format(timestampFormat))
}

// Summary returns a formatted summary string of the environment report
func (r *EnvironmentReport) Summary() string {
	const summaryTemplate = `--- SCMPUFF DEBUG DUMPTOOL ---
User:          %s
Timestamp:     %s
git version:   %s
os.Getwd:      %q
show-toplevel: %q
show-cdup:     %q

$ git status --branch --porcelain=v1
%s`
	timestamp := r.timestamp.Format(timestampFormat)
	return fmt.Sprintf(summaryTemplate, r.user, timestamp, r.gitVersion, r.cwd, r.gitRoot, r.cdupPath, r.statusData)
}

func writeSummaryFile(outputDir string, samplePrefix string, content string, verbose bool) error {
	outputDirRoot, err := os.OpenRoot(outputDir)
	if err != nil {
		return fmt.Errorf("could not open %q root: %w", outputDir, err)
	}
	defer outputDirRoot.Close()

	filename := samplePrefix + "-summary.txt"
	f, err := outputDirRoot.Create(filename)
	if err != nil {
		return fmt.Errorf("failed creating summary file: %w", err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed writing summary file: %w", err)
	}
	if verbose {
		log.Printf("SUCCESS: summary written to %s", filename)
	}
	return nil
}

// createPorcelainSamples generates all the required git status porcelain format
// files and writes them to the specified output directory, using samplePrefix
// in the filename template.
func createPorcelainSamples(outputDir string, samplePrefix string, verbose bool) error {
	samples := []struct {
		args   []string
		suffix string
	}{
		{[]string{"status", "--branch", "--porcelain=v1"}, ".porcelain.v1.txt"},        //   git status --branch --porcelain=v1    > $TMPDIR/$NAME-porcelain.v1.txt
		{[]string{"status", "--branch", "--porcelain=v1", "-z"}, ".porcelain.v1z.bin"}, //   git status --branch --porcelain=v1 -z > $TMPDIR/$NAME-porcelain.v1z.bin
		{[]string{"status", "--branch", "--porcelain=v2"}, ".porcelain.v2.txt"},        //   git status --branch --porcelain=v2    > $TMPDIR/$NAME-porcelain.v2.txt
		{[]string{"status", "--branch", "--porcelain=v2", "-z"}, ".porcelain.v2z.bin"}, //   git status --branch --porcelain=v2 -z > $TMPDIR/$NAME-porcelain.v2z.bin
	}

	// use new os.Root in go1.24 to prevent traversal outside of outputDir for
	// safety, probably not actually needed here, but good practice to use most
	// restrictive defaults.
	outputDirRoot, err := os.OpenRoot(outputDir)
	if err != nil {
		return fmt.Errorf("could not open %q root: %w", outputDir, err)
	}
	defer outputDirRoot.Close()

	for _, sample := range samples {
		filename := samplePrefix + sample.suffix
		f, err := outputDirRoot.Create(filename)
		if err != nil {
			return fmt.Errorf("failed creating output file %q: %w", filename, err)
		}
		defer f.Close()

		cmd := exec.Command("git", sample.args...)
		cmd.Stdout = f
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error running git command %v: %w", sample.args, err)
		}

		if verbose {
			log.Printf("SUCCESS: git %v > %s", sample.args, filename)
		}
	}

	return nil
}

// createArchive creates a ZIP archive containing the entire contents of the tmpdir.
// it is created in the current working directory for user convenience.
func createArchive(tmpdir string, slug string) error {
	tmpdirRoot, err := os.OpenRoot(tmpdir)
	if err != nil {
		return fmt.Errorf("could not open %q root: %w", tmpdir, err)
	}
	defer tmpdirRoot.Close()

	archiveName := slug + ".zip"
	archiveFile, err := os.Create(archiveName)
	if err != nil {
		return fmt.Errorf("failed to create archive file %q: %w", archiveName, err)
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	if err := zipWriter.AddFS(tmpdirRoot.FS()); err != nil {
		return fmt.Errorf("failed to add files to archive %q: %w", archiveName, err)
	}

	log.Printf("Created archive %s", archiveName)
	return nil
}
