package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

//GetContentLength
func GetContentLength(url string) (err error, contentLength uint64) {
	res, err := http.Head(url)
	if err != nil {
		return err, 0
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("http error: status code %d", res.StatusCode), 0
	}
	defer res.Body.Close()

	length, err := strconv.ParseUint(res.Header.Get("Content-Length"), 0, 64)
	if err != nil {
		return err, 0
	}

	return nil, length
}

//RangeDownload
func RangeDownload(url string, startPos int64, rangeByte int64) error {
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
