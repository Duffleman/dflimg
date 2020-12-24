package cmd

import (
	"context"
	"time"

	"dflimg"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RemoveShortcutCmd = &cobra.Command{
	Use:     "remove-shortcut {query} {shortcut}",
	Aliases: []string{"rsc"},
	Short:   "Remove a shortcut",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		query := args[0]
		shortcut := args[1]

		err := removeShortcut(ctx, query, shortcut)
		if err != nil {
			return err
		}

		duration := time.Now().Sub(startTime)

		log.Infof("Done in %s", duration)

		return nil
	},
}

func removeShortcut(ctx context.Context, query, shortcut string) error {
	body := &dflimg.ChangeShortcutRequest{
		IdentifyResource: dflimg.IdentifyResource{
			Query: query,
		},
		Shortcut: shortcut,
	}

	return makeClient().RemoveShortcut(ctx, body)
}
