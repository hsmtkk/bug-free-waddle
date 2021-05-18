package login

import (
	"log"

	"github.com/hsmtkk/bug-free-waddle/command/env"
	"github.com/hsmtkk/bug-free-waddle/http"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:  "login",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	username, password, err := env.UsernamePassword()
	if err != nil {
		log.Fatal(err)
	}
	getter, err := http.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := getter.Login(username, password); err != nil {
		log.Fatal(err)
	}
}
