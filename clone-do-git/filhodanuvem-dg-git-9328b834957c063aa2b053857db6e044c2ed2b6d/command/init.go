package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/devgymbr/fit/index"
	"github.com/spf13/cobra"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "start a new fit repository",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(".fit")
		if !os.IsNotExist(err) {
			log.Fatalf("failed to initialize fit repository: .fit folder already exists")
		}

		if err := os.Mkdir(".fit", 0755); err != nil {
			log.Fatalf("failed to initialize fit repository: %s", err)
		}

		err = func() error {
			if err := os.Mkdir(filepath.Join(".", ".fit", "objects"), 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(filepath.Join(".", ".fit", "index"), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}

			defer file.Close()

			head, err := os.OpenFile(filepath.Join(".", ".fit", "HEAD"), os.O_CREATE, 0644)
			if err != nil {
				return err
			}

			defer head.Close()

			var idx index.Index
			idx.Objects = make(map[string]index.Entry)
			content, err := json.Marshal(idx)
			if err != nil {
				return err
			}
			_, err = file.Write(content)

			return err
		}()

		if err != nil {
			err2 := os.RemoveAll(".fit")
			if err2 != nil {
				log.Fatalf("failed to initialize fit repository: corrupted .fit folder: %s", err)
			}
			log.Fatalf("failed to initialize fit repository: %s", err)
		}

		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to initialize fit repository: %s", err)
		}
		fmt.Printf("Initialized empty fit repository in %s \n", filepath.Join(pwd, ".fit"))
	},
}
