package command

import (
	"fmt"
	"log"
	"strconv"

	"github.com/devgymbr/fit/commit"
	"github.com/spf13/cobra"
)

var Log = &cobra.Command{
	Use:   "log",
	Short: "list commits",
	Run: func(cmd *cobra.Command, args []string) {
		current, err := commit.GetHEAD()
		if err != nil {
			log.Fatalf("not possible to get HEAD: %s", err)
		}

		parent, err := commit.GetByHash(current.ParentHash)
		if err != nil {
			log.Fatalf("not possible to get HEAD: %s", err)
		}

		limit := int64(5)
		if len(args) > 0 {
			value, err := strconv.ParseInt(args[0], 10, 64)
			if err == nil {
				limit = value
			}
		}

		for i := int64(0); i < limit; i++ {
			current.Parent = &parent
			fmt.Printf("%s\n", current)
			if current.ParentHash == "" {
				return
			}

			current = parent
			parent, err = commit.GetByHash(parent.ParentHash)
			if err != nil {
				log.Fatalf("not possible to get commit %s: %s", parent.ParentHash, err)
			}
		}
	},
}
