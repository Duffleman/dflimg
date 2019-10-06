package cmd

import (
	b "github.com/gen2brain/beeep"
	log "github.com/sirupsen/logrus"
)

// AppName for notifications
const AppName = "DFL IMG"

func notify(title, body string) {
	err := b.Notify(title, body, "")
	if err != nil {
		log.Warn(err)
	}
}
