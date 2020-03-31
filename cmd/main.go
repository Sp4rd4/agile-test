package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sp4rd4/fuzzyelem"
	"github.com/urfave/cli/v2"
)

func main() {
	var id string

	app := &cli.App{
		Name:  "fuzzyelem",
		Usage: "search in a target document for htlm element similar to element from source document",

		UsageText: "fuzzyelem [GLOBAL OPTIONS] source target",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "id",
				Value:       "make-everything-ok-button",
				Usage:       "source document element 'id' attribute value",
				Destination: &id,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 2 {
				return cli.Exit("not enough arguments", 1)
			}
			if c.Args().Len() > 2 {
				return cli.Exit("too much arguments", 1)
			}
			path, err := fuzzyelem.Search(id, c.Args().Get(0), c.Args().Get(1))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
			fmt.Println(path)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
