package command

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/devgymbr/fit/commit"
	"github.com/devgymbr/fit/index"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var idx index.Index
var head commit.Entry
var notStagedFiles []string
var stagedFiles []string
var untrackedFiles []string

var Status = &cobra.Command{
	Use:   "status",
	Short: "check the status of the repository",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		idx, err = index.Build()
		if err != nil {
			log.Fatalf("could not build index: %v", err)
		}

		head, _ = commit.GetHEAD()

		filepath.Walk(".", classifyFiles)

		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		if len(stagedFiles) > 0 {
			fmt.Println("Changes on stage:")
			fmt.Println("  (no support to \"fit rm <file>\" yet)")
			fmt.Println("    " + green(strings.Join(stagedFiles, "\n    ")))
		}

		if len(notStagedFiles) > 0 {
			fmt.Println("Changes not staged for commit:")
			fmt.Println("  (use \"fit add <file>\" to update what will be committed)")
			fmt.Println("    " + red(strings.Join(notStagedFiles, "\n    ")))
		}

		if len(untrackedFiles) > 0 {
			fmt.Println("untracked files:")
			fmt.Println("  (use \"fit add <file>\" to update what will be committed)")
			fmt.Println("    " + red(strings.Join(untrackedFiles, "\n    ")))
		}

		if len(stagedFiles) == 0 && len(notStagedFiles) == 0 && len(untrackedFiles) == 0 {
			fmt.Println("nothing to commit, working tree clean")
		}
	},
}

func classifyFiles(path string, fileInfo fs.FileInfo, err error) error {
	if fileInfo.IsDir() && (fileInfo.Name() == ".fit" || fileInfo.Name() == ".git") {
		return filepath.SkipDir
	}

	if fileInfo.IsDir() {
		return nil
	}

	if fileInfo.Name() == "." || fileInfo.Name() == ".." {
		return nil
	}

	if _, ok := idx.Objects[path]; !ok {
		untrackedFiles = append(untrackedFiles, path)
		return nil
	}

	if idx.Objects[path].LastModified.Before(fileInfo.ModTime()) {
		notStagedFiles = append(notStagedFiles, path)
		return nil
	}

	if head.Index.Objects[path].Hash != idx.Objects[path].Hash {
		stagedFiles = append(stagedFiles, path)
		return nil
	}

	return nil
}
