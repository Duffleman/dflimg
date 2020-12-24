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
	rootCmd.AddCommand(cli.AddShortcutCmd)
	rootCmd.AddCommand(cli.CopyURLCmd)
	rootCmd.AddCommand(cli.DeleteResourceCmd)
	rootCmd.AddCommand(cli.RemoveShortcutCmd)
	rootCmd.AddCommand(cli.ScreenshotCmd)
	rootCmd.AddCommand(cli.SetNSFWCmd)
	rootCmd.AddCommand(cli.ShortenURLCmd)
	rootCmd.AddCommand(cli.UploadSignedCmd)
	rootCmd.AddCommand(cli.ViewDetailsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "dflimg",
	Short: "CLI tool to upload images to a dflimg server",
	Long:  "A CLI tool to manage files and URLs being uploaded and removed from your chosen dflimg server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
