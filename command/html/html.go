package html

import (
	"fmt"
	"log"

	"github.com/hsmtkk/bug-free-waddle/http"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:  "html url",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func run(url string) {
	getter, err := http.New()
	if err != nil {
		log.Fatal(err)
	}
	html, err := getter.GetAlbum(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html)
}
