package resolve

import (
	"fmt"
	"log"

	"github.com/hsmtkk/bug-free-waddle/command/file"
	"github.com/hsmtkk/bug-free-waddle/http"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:  "resolve datasrc.txt",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func run(datasrc string) {
	paths, err := file.ReadLines(datasrc)
	if err != nil {
		log.Fatal(err)
	}
	getter, err := http.New()
	if err != nil {
		log.Fatal(err)
	}
	urls, err := getter.ResolveDataSrc(paths)
	if err != nil {
		log.Fatal(err)
	}
	for _, url := range urls {
		fmt.Println(url)
	}
}
