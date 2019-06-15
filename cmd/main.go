package main

import (
	"goDownloader/downloader"
	"log"
	"os"

	"github.com/urfave/cli"
)

func defineApp(app *cli.App) {
	app.Name = "goDownloader"
	app.Usage = "file downloader made of Go."
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "pallarel, p",
			Value: 0,
			Usage: "pallarel number to download file",
		},
		cli.StringFlag{
			Name:  "url, u",
			Usage: "file url to download",
		},
	}
	app.Action = func(c *cli.Context) error {
		downloader.Download("WIP")
		return nil
	}

}

func main() {
	app := cli.NewApp()
	defineApp(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
