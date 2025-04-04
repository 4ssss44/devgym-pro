package command

import (
	"fmt"
	"log"

	"github.com/devgymbr/fit/commit"
	"github.com/devgymbr/fit/objects"
	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var Show = &cobra.Command{
	Use:   "show",
	Short: "show all commit changes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("commit hash is required")
		}

		current, err := commit.GetByHash(args[0])
		if err != nil {
			log.Fatalf("not possible to get commit %s: %s", args[0], err)
		}

		if current.Hash == "" {
			log.Fatalf("commit %s not found", args[0])
		}

		parent, err := commit.GetByHash(current.ParentHash)
		if err != nil {
			log.Fatalf("not possible to get commit %s: %s", args[0], err)
		}

		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		for _, file := range current.Index.Objects {
			if _, exists := parent.Index.Objects[file.Path]; !exists {
				fmt.Printf("%s:\n", yellow(file.Path))
				content, err := objects.ReadBlob(file.Hash)
				if err != nil {
					log.Fatalf("not possible to read blob %s from commit %s: %s", file.Hash, current.Hash, err)
				}

				fmt.Printf("%s\n", green(string(content)))
				continue
			}

			if parent.Index.Objects[file.Path].Hash != file.Hash {
				fmt.Printf("%s:\n", yellow(file.Path))
				currentContent, err := objects.ReadBlob(file.Hash)
				if err != nil {
					log.Fatalf("not possible to read blob %s from parent commit %s: %s", file.Hash, parent.Hash, err)
				}

				parentContent, err := objects.ReadBlob(parent.Index.Objects[file.Path].Hash)
				if err != nil {
					log.Fatalf("not possible to read blob %s: %s", file.Hash, err)
				}

				diff := diffmatchpatch.New()
				contentDiff := diff.DiffMain(string(parentContent), string(currentContent), false)

				fmt.Printf("%s\n", diff.DiffPrettyText(contentDiff))
			}
		}

	},
}
