package command

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/devgymbr/fit/commit"
	"github.com/devgymbr/fit/index"
	"github.com/spf13/cobra"
)

var Commit = &cobra.Command{
	Use:   "commit",
	Short: "create commit based on tracked files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("commit message is required")
		}

		idx, err := index.Build()
		if err != nil {
			log.Fatalf("could not build index: %v", err)
		}

		authorName := getEnvOrDefault("FIT_AUTHOR_NAME", "Rob Pike")
		authorEmail := getEnvOrDefault("FIT_AUTHOR_EMAIL", "guest@gopher.com")
		c := commit.Entry{
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			AuthorDate:  time.Now(),
			Message:     args[0],
			Index:       idx,
		}

		head, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("could not get HEAD: %v", err)
		}

		if head.Hash != "" {
			c.ParentHash = head.Hash
		}

		hash, err := commit.Write(c)
		if err != nil {
			log.Fatalf("could not write commit: %v", err)
		}
		fmt.Printf("[%s] %s\n", hash[0:6], c.Message)
	},
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
