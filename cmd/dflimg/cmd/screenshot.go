package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/cuvva/ksuid"
	"github.com/spf13/cobra"
)

const screenshotCmd = "screencapture -i"

var ScreenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Take a screenshot & upload it",
	Long:  "Take a screenshot and upload it to a DFLIMG server",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tmpDir := fmt.Sprintf("/tmp/%s.png", ksuid.Generate("file").String())

		err := exec.CommandContext(ctx, "screencapture", "-i", tmpDir).Run()
		if err != nil {
			return err
		}
		defer os.Remove(tmpDir)

		_, err = os.Stat(tmpDir)
		if os.IsNotExist(err) {
			return nil
		}

		return UploadSignedCmd.RunE(cmd, []string{tmpDir})
	},
}
