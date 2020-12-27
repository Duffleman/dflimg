package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/cuvva/ksuid-go"
	"github.com/spf13/cobra"
)

const screenshotCmd = "screencapture -i"
const timeout = 1 * time.Minute

var ScreenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Take a screenshot & upload it",
	Long:  "Take a screenshot and upload it to a DFLIMG server",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		tmpName := fmt.Sprintf("%s-*.png", ksuid.Generate("file").String())
		out, err := ioutil.TempFile("", tmpName)
		if err != nil {
			return err
		}
		defer out.Close()

		err = exec.CommandContext(ctx, "screencapture", "-i", out.Name()).Run()
		if err != nil {
			return err
		}
		defer os.Remove(out.Name())

		tmpFile, err := os.Stat(out.Name())
		if os.IsNotExist(err) {
			return nil
		}

		if tmpFile.Size() == 0 {
			notify("Cancelled", "No image was captured.")
			return nil
		}

		return UploadSignedCmd.RunE(cmd, []string{out.Name()})
	},
}
