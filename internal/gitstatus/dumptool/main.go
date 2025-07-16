// Dumptool is a command-line tool for dumping git status information
// in a machine-readable format for debugging or analysis purposes.
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

var (
	wantSummary   = flag.Bool("summary", true, "summary of status for debugging")
	wantPorcelain = flag.Bool("porcelain", false, "keep porcelain files for test or debug purposes")
	wantArchive   = flag.Bool("archive", false, "create compressed archive of data")
	verbose       = flag.Bool("v", false, "enable verbose output for debugging")
)

func checkFatalErr(msg string, err error) {
	if err != nil {
		log.Fatalf("fatal: %s: %v", msg, err)
	}
}

func username() string {
	user, err := user.Current()
	if err != nil {
		return "anonymous"
	}
	return user.Username
}

const summaryTemplate = `--- SCMPUFF DEBUG DUMPTOOL ---
User:          %s
Timestamp:     %s
os.Getwd:      %q
show-toplevel: %q
show-cdup:     %q

$ git status --branch --porcelain=v1
%s`

func main() {
	flag.Parse()

	wd, err := os.Getwd()
	checkFatalErr("failed to retrieve current working directory", err)

	// first we need to determine the git project root, which we also use for
	// naming our output files.  if we can't get that, we won't be able to
	// proceed anyhow since we're likely not in a git repository, so die.
	toplevelBytes, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	checkFatalErr("failed to determine git project root", err)

	// some metadata for attribution and filenames
	user := username()
	basename := filepath.Base(string(bytes.TrimSpace(toplevelBytes))) // base name of the repo
	timestamp := time.Now().Format("20060102150405")
	slug := fmt.Sprintf("%s-%s-%s", user, basename, timestamp)

	var summary string
	if *wantSummary {
		cdupBytes, err := exec.Command("git", "rev-parse", "--show-cdup").Output()
		checkFatalErr("failed on call to --show-cdup", err)

		statusBytes, err := exec.Command("git", "status", "--branch", "--porcelain=v1").Output()
		checkFatalErr("git status error", err)

		summary = fmt.Sprintf(summaryTemplate, user, timestamp, wd, bytes.TrimSpace(toplevelBytes), bytes.TrimSpace(cdupBytes), statusBytes)
		fmt.Println(summary)
	}

	if *wantPorcelain || *wantArchive { // archive requires porcelain files
		// create a temporary directory for porcelain files.
		// if user *specifically* asked for porcelain files, keep them around afterwards,
		// otherwise clean up tmpdir after we are done.
		tmpdir, err := os.MkdirTemp("", "scmpuff-dumptool-")
		if err != nil {
			log.Fatalf("fatal: could not make tmpdir: %v", err)
		}
		if !*wantPorcelain {
			defer func() {
				err := os.RemoveAll(tmpdir)
				if *verbose {
					log.Printf("VERBOSE: removed tmpdir %s, error = %v", tmpdir, err)
				}
			}()
		}

		// write porcelain files to tmpdir
		if err := writePorcelainSamples(tmpdir, slug, summary); err != nil {
			log.Fatalf("fatal: failed to write porcelain samples: %v", err)
		}

		// if porcelain files were specifically requested, let user know
		if *wantPorcelain {
			log.Printf("Porcelain files saved to: %s", tmpdir)
		}

		if *wantArchive {
			createArchive(tmpdir, slug)
		}
	}

}

// writePorcelainSamples generates all the required git status porcelain format
// files and writes them to the specified output directory, using samplePrefix
// in the filename template.
//
// if includeSummary is not empty, a summary.txt file will also be included
// with the contents of includeSummary.
func writePorcelainSamples(outputDir string, samplePrefix string, includeSummary string) error {
	samples := []struct {
		args   []string
		suffix string
	}{
		{[]string{"status", "--branch", "--porcelain=v1"}, ".porcelain-v1.txt"},        //   git status --branch --porcelain=v1    > $TMPDIR/$NAME-porcelain-v1.txt
		{[]string{"status", "--branch", "--porcelain=v1", "-z"}, ".porcelain-v1z.txt"}, //   git status --branch --porcelain=v1 -z > $TMPDIR/$NAME-porcelain-v1z.txt
		{[]string{"status", "--branch", "--porcelain=v2"}, ".porcelain-v2.txt"},        //   git status --branch --porcelain=v2    > $TMPDIR/$NAME-porcelain-v2.txt
		{[]string{"status", "--branch", "--porcelain=v2", "-z"}, ".porcelain-v2z.txt"}, //   git status --branch --porcelain=v2 -z > $TMPDIR/$NAME-porcelain-v2z.txt
	}

	// use new os.Root in go1.24 to prevent traversal outside of outputDir for
	// safety, probably not actually needed here, but good practice to use most
	// restrictive defaults.
	outputDirRoot, err := os.OpenRoot(outputDir)
	if err != nil {
		return fmt.Errorf("could not open %q root: %v", outputDir, err)
	}
	defer outputDirRoot.Close()

	for _, sample := range samples {
		filename := samplePrefix + sample.suffix
		f, err := outputDirRoot.Create(filename)
		if err != nil {
			return fmt.Errorf("failed creating output file %q: %v", filename, err)
		}
		defer f.Close()

		cmd := exec.Command("git", sample.args...)
		cmd.Stdout = f
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error running git command %v: %v", sample.args, err)
		}

		if *verbose {
			log.Printf("SUCCESS: git %v > %s", sample.args, filename)
		}
	}

	if includeSummary != "" {
		filename := samplePrefix + "-summary.txt"
		f, err := outputDirRoot.Create(filename)
		if err != nil {
			return fmt.Errorf("failed creating summary file: %v", err)
		}
		defer f.Close()
		if _, err := f.WriteString(includeSummary); err != nil {
			return fmt.Errorf("failed writing summary file: %v", err)
		}
		if *verbose {
			log.Printf("SUCCESS: summary written to %s", filename)
		}
	}

	return nil
}

// createArchive creates a ZIP archive containing the entire contents of the tmpdir.
// it is created in the current working directory for user convenience.
func createArchive(tmpdir string, slug string) error {
	tmpdirRoot, err := os.OpenRoot(tmpdir)
	if err != nil {
		return fmt.Errorf("could not open %q root: %v", tmpdir, err)
	}
	defer tmpdirRoot.Close()

	archiveName := slug + ".zip"
	archiveFile, err := os.Create(archiveName)
	if err != nil {
		return fmt.Errorf("failed to create archive file %q: %v", archiveName, err)
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	if err := zipWriter.AddFS(tmpdirRoot.FS()); err != nil {
		return fmt.Errorf("failed to add files to archive %q: %v", archiveName, err)
	}

	log.Printf("Created archive %s", archiveName)
	return nil
}
