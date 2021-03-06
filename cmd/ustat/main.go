package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

const version = "0.2.0"

func main() {
	app := cli.NewApp()
	app.Name = "ustat"
	app.Version = version
	app.Usage = "Unified system statistics collector"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Pekka Enberg",
			Email: "penberg@iki.fi",
		},
	}
	app.HideHelp = true
	app.Commands = []cli.Command{
		recordCommand,
		reportCommand,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
