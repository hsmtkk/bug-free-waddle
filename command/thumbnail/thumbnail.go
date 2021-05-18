package thumbnail

import (
	"log"

	"github.com/hsmtkk/bug-free-waddle/command/file"
	"github.com/hsmtkk/bug-free-waddle/http"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:  "thumbnail thumbnail.txt",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func run(path string) {
	urls, err := file.ReadLines(path)
	if err != nil {
		log.Fatal(err)
	}
	getter, err := http.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := getter.GetThumbnail(urls); err != nil {
		log.Fatal(err)
	}
}
