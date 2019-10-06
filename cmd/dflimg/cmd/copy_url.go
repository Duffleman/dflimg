package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cuvva/ksuid"
	"github.com/spf13/cobra"
)

var CopyURLCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"c"},
	Short:   "Copy from a URL",
	Long:    "Copy from a URL to the dflimg server",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, err := downloadFile(args[0])
		if err != nil {
			return err
		}
		defer os.Remove(*filePath)

		return UploadSignedCmd.RunE(cmd, []string{*filePath})
	},
}

func downloadFile(urlStr string) (*string, error) {
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
