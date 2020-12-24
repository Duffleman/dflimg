package cmd

import (
	"context"
	"fmt"
	"time"

	"dflimg"

	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var AddShortcutCmd = &cobra.Command{
	Use:     "add-shortcut {query} {shortcut}",
	Aliases: []string{"asc"},
	Short:   "Add a shortcut",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		query := args[0]
		shortcut := args[1]

		err := addShortcut(ctx, query, shortcut)
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(fmt.Sprintf("%s/:%s", rootURL(), shortcut))
		if err != nil {
			log.Warn("Could not copy to clipboard.")
		}

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s", duration)

		return nil
	},
}

func addShortcut(ctx context.Context, query, shortcut string) error {
	body := &dflimg.ChangeShortcutRequest{
		IdentifyResource: dflimg.IdentifyResource{
			Query: query,
		},
		Shortcut: shortcut,
	}

	return makeClient().AddShortcut(ctx, body)
}
