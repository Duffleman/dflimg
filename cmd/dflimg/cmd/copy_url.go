package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cuvva/ksuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CopyURLCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"c"},
	Short:   "Copy from a URL",
	Long:    "Copy from a URL to the dflimg server",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		startTime := time.Now()

		urlStr := args[0]
		shortcuts := cmd.Flag("shortcuts")
		nsfw := cmd.Flag("nsfw")

		rootURL := viper.Get("ROOT_URL").(string)
		authToken := viper.Get("AUTH_TOKEN").(string)

		filePath, err := downloadFile(rootURL, authToken, urlStr)
		if err != nil {
			return err
		}
		defer os.Remove(*filePath)

		body, err := sendFile(rootURL, authToken, *filePath, shortcuts, nsfw)
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(body.URL)
		if err != nil {
			fmt.Println("Could not copy to clipboard. Please copy the URL manually")
		}

		duration := time.Now().Sub(startTime)

		fmt.Printf("Done in %s: %s\n", duration, body.URL)

		return nil
	},
}

func downloadFile(rootURL, authToken, urlStr string) (*string, error) {
	fileToCopyRes, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer fileToCopyRes.Body.Close()

	tmpName := ksuid.Generate("tmpfile").String()
	tmpPath := fmt.Sprintf("/tmp/%s", tmpName)

	out, err := os.Create(tmpPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, fileToCopyRes.Body)
	if err != nil {
		return nil, err
	}

	return &tmpPath, nil
}
