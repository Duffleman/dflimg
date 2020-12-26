package cmd

import (
	"dflimg"

	"github.com/atotto/clipboard"
	b "github.com/gen2brain/beeep"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// AppName for notifications
const AppName = "DFL IMG"

func notify(title, body string) {
	err := b.Notify(title, body, "")
	if err != nil {
		log.Warn(err)
	}
}

func makeClient() dflimg.Service {
	authToken := viper.Get("AUTH_TOKEN").(string)

	return dflimg.NewClient(rootURL(), authToken)
}

func rootURL() string {
	return viper.Get("ROOT_URL").(string)
}

func writeClipboard(in string) {
	err := clipboard.WriteAll(in)
	if err != nil {
		log.Warn("Could not copy to clipboard.")
	}
}
