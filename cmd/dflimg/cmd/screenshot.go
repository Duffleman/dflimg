package cmd

import (
	"context"
	"io/ioutil"
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

		tmpName := ksuid.Generate("file").String()
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

		_, err = os.Stat(out.Name())
		if os.IsNotExist(err) {
			return nil
		}

		return UploadSignedCmd.RunE(cmd, []string{out.Name()})
	},
}
