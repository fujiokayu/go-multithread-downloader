package main

import (
	"goDownloader/downloader"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	err := cli.NewApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	downloader.Download("WIP")
}
