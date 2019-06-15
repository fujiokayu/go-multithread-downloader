package downloader

import (
	"io"
	"net/http"
	"os"
	"path"
)

//Download
func Download(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Create the file
	out, err := os.Create(path.Base(url))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)
	return err
}
