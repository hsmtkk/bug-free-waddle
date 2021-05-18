package main

import (
	"log"

	"github.com/hsmtkk/bug-free-waddle/command/datasrc"
	"github.com/hsmtkk/bug-free-waddle/command/html"
	"github.com/hsmtkk/bug-free-waddle/command/login"
	"github.com/hsmtkk/bug-free-waddle/command/resolve"
	"github.com/hsmtkk/bug-free-waddle/command/thumbnail"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use: "en-photo-screenshot",
}

func init() {
	command.AddCommand(html.Command)
	command.AddCommand(datasrc.Command)
	command.AddCommand(resolve.Command)
	command.AddCommand(thumbnail.Command)
	command.AddCommand(login.Command)
}

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
