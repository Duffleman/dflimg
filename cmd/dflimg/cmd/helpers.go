package cmd

import (
	"dflimg"

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
	rootURL := viper.Get("ROOT_URL").(string)
	authToken := viper.Get("AUTH_TOKEN").(string)

	return dflimg.NewClient(rootURL, authToken)
}

func rootURL() string {
	return viper.Get("ROOT_URL").(string)
}
