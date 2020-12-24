package cmd

import (
	"dflimg"
	dhttp "dflimg/cmd/dflimg/http"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var ViewDetailsCmd = &cobra.Command{
	Use:     "view {query}",
	Aliases: []string{"v"},
	Short:   "View details of a resource",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		rootURL, authToken := setup()
		c := dhttp.New(rootURL, authToken)

		res := &dflimg.Resource{}
		req := &dflimg.IdentifyResource{
			Query: query,
		}

		err := c.JSONRequest("POST", "view_details", req, res)
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
