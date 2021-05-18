package datasrc

import (
	"fmt"
	"log"
	"os"

	"github.com/hsmtkk/bug-free-waddle/datasrc"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:  "datasrc album.html",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func run(htmlPath string) {
	content, err := os.ReadFile(htmlPath)
	if err != nil {
		log.Fatalf("failed to open HTML file; %s; %s", htmlPath, err)
	}
	urls, err := datasrc.New().SelectDataSrc(string(content))
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range urls {
		fmt.Println(u)
	}
}
