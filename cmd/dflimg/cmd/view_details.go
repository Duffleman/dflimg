package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"dflimg"

	"github.com/spf13/cobra"
)

var ViewDetailsCmd = &cobra.Command{
	Use:     "view {query}",
	Aliases: []string{"v"},
	Short:   "View details of a resource",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		query := args[0]

		res, err := makeClient().ViewDetails(ctx, &dflimg.IdentifyResource{
			Query: query,
		})
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(b))

		return nil
	},
}
