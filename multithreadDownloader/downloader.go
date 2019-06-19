package multithreadDownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
)

type DownlodeClient struct {
	URL             string
	ContentLength   int64
	HasAcceptRanges bool
	ThreadNumber    int
	IsReady         bool
}

func (DownlodeClient *DownlodeClient) SetResponceHeader() error {
	fmt.Println("setResponceHeader of ", DownlodeClient.URL)
	res, err := http.Head(DownlodeClient.URL)
	if err != nil {
		return fmt.Errorf("failed to get Header: %s", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("http error: status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	DownlodeClient.ContentLength, err = strconv.ParseInt(res.Header.Get("Content-Length"), 0, 64)
	if err != nil {
		return fmt.Errorf("failed to get content-length: %s", err)
	}

	DownlodeClient.HasAcceptRanges = (res.Header.Get("Accept-Ranges") == "bytes")
	DownlodeClient.IsReady = DownlodeClient.HasAcceptRanges

	return nil
}

//Download startPos int64, rangeByte int64
func (downlodeClient DownlodeClient) Download(threadNumber int) error {

	if threadNumber == 0 {
		threadNumber = runtime.NumCPU()
	}

	if !downlodeClient.IsReady {
		return fmt.Errorf("DownlodeClient is not ready")
	}

	fmt.Sprintf("download %d parallels", threadNumber)

	payloadSize := downlodeClient.ContentLength / int64(threadNumber)
	remaindSize := downlodeClient.ContentLength
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the file
	out, err := os.Create(path.Base(downlodeClient.URL))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()

	for i := 0; remaindSize > 0; i++ {
		fmt.Println("Go: ", remaindSize)
		startPos := downlodeClient.ContentLength - remaindSize
		endPos := startPos + payloadSize
		if endPos > downlodeClient.ContentLength {
			endPos = downlodeClient.ContentLength
		}
		remaindSize -= payloadSize

		req, err := http.NewRequest("GET", downlodeClient.URL, nil)
		if err != nil {
			//				return err
		}
		req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", startPos, endPos))
		req = req.WithContext(ctx)
		fmt.Println(i, ": done")

		var client http.Client
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Do i: ", err)
		}
		defer res.Body.Close()
		//var buf bytes.Buffer
		_, err = io.Copy(out, res.Body)
		if err != nil {
			fmt.Println("Copy i: ", err)
		}

		// Write the body to file
		_, err = io.Copy(out, res.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

	}

	// Write the body to file
	//_, err = io.Copy(out, res.Body)
	fmt.Println("done")
	return err
}
