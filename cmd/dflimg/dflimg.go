package main

import (
	"fmt"
	"os"

	cli "dflimg/cmd/dflimg/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Load env variables
	viper.SetEnvPrefix("DFLIMG")
	viper.SetDefault("ROOT_URL", "https://dfl.mn")

	viper.AutomaticEnv()

	// Register commands
	rootCmd.AddCommand(cli.UploadCmd)

	// handle command argumetns
	cli.UploadCmd.Flags().StringP("shortcuts", "s", "", "A CSV of shortcuts to apply to the uploaded file")
	cli.UploadCmd.Flags().BoolP("nsfw", "n", false, "Is the resource NSFW?")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "dflimg",
	Short: "CLI tool to upload images to a dflimg server",
	Long:  "A CLI tool to manage files being uploaded, labeled, and removed from your chosen dflimg server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
