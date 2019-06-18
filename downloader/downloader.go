package downloader

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
)

type DownlodeClient struct {
	url            string
	contentType    string
	contentLength  uint64
	responseHeader *http.Response
}

func (DownlodeClient *DownlodeClient) setResponceHeader() error {
	fmt.Println("setResponceHeader of ", DownlodeClient.url)
	res, err := http.Head(DownlodeClient.url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("http error: status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	DownlodeClient.responseHeader = res
	return nil
}

func HasAcceptRanges(dc DownlodeClient) (error, bool) {
	err := dc.setResponceHeader()
	if err != nil {
		return err, false
	}
	res := dc.responseHeader.Header.Get("Accept-Ranges")

	return nil, res == "bytes"
}

//GetContentLength
func GetContentLength(url string) (error, uint64) {
	dc := &DownlodeClient{url, "", 0, nil}
	err := dc.setResponceHeader()
	if err != nil {
		return err, 0
	}
	fmt.Println(HasAcceptRanges(*dc))
	length, err := strconv.ParseUint(dc.responseHeader.Header.Get("Content-Length"), 0, 64)
	if err != nil {
		return err, 0
	}

	return nil, length
}

//RangeDownload
func RangeDownload(url string, startPos int64, rangeByte int64) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//client := &http.Client{Timeout: time.Duration(10) * time.Second}

	defer req.Body.Close()

	// Create the file
	out, err := os.Create(path.Base(url))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	//_, err = io.Copy(out, res.Body)
	return err
}
