package cmd

import (
	"context"
	"fmt"
	"time"

	"dflimg"
	"dflimg/lib/cher"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var AddShortcutCmd = &cobra.Command{
	Use:     "add-shortcut {query} {shortcut}",
	Aliases: []string{"asc"},
	Short:   "Add a shortcut",
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

		err = addShortcut(ctx, query, shortcut)
		if err != nil {
			return err
		}

		writeClipboard(fmt.Sprintf("%s/:%s", rootURL(), shortcut))
		notify("Added shortcut", fmt.Sprintf("%s/:%s", rootURL(), shortcut))

		log.Infof("Done in %s", time.Now().Sub(startTime))

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

func handleShortcutInput(args []string) (string, string, error) {
	if len(args) == 2 {
		return args[0], args[1], nil
	}

	queryPrompt := promptui.Prompt{
		Label: "Query",
		Validate: func(input string) error {
			if len(input) >= 1 {
				return nil
			}

			return cher.New("missing_query", nil)
		},
	}

	shortcutPrompt := promptui.Prompt{
		Label: "Shortcut",
		Validate: func(input string) error {
			if len(input) >= 1 {
				return nil
			}

			return cher.New("missing_shortcut", nil)
		},
	}

	query, err := queryPrompt.Run()
	if err != nil {
		return "", "", err
	}

	shortcut, err := shortcutPrompt.Run()
	if err != nil {
		return "", "", err
	}

	return query, shortcut, nil
}
