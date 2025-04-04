package command

import (
	"log"
	"os"

	"github.com/devgymbr/fit/index"
	"github.com/devgymbr/fit/objects"
	"github.com/spf13/cobra"
)

var Add = &cobra.Command{
	Use:   "add",
	Short: "add file to fit staging area",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("missing file argument")
		}

		for _, path := range args {
			file, err := os.OpenFile(path, os.O_RDONLY, 0655)
			if err != nil {
				log.Fatalf("cannot open file %s: %s", path, err)
			}
			defer file.Close()

			blob, err := objects.CreateBlob(file)
			if err != nil {
				log.Fatalf("cannot create blob: %s", err)
			}
			defer blob.Close()

			err = index.AddBlob(file, blob)
			if err != nil {
				log.Fatalf("cannot add blob: %s", err)
			}
		}
	},
}
