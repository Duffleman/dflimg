package cmd

import (
	"context"
	"time"

	"dflimg"
	"dflimg/lib/cher"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RemoveShortcutCmd = &cobra.Command{
	Use:     "remove-shortcut {query} {shortcut}",
	Aliases: []string{"rsc"},
	Short:   "Remove a shortcut",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 2 || len(args) == 0 {
			return nil
		}

		return cher.New("missing_arguments", nil)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		startTime := time.Now()

		query, shortcut, err := handleShortcutInput(args)
		if err != nil {
			return err
		}

		err = removeShortcut(ctx, query, shortcut)
		if err != nil {
			return err
		}

		notify("Removed shortcut", shortcut)

		log.Infof("Done in %s", time.Now().Sub(startTime))

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
