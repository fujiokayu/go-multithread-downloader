package multithreadDownloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"

	"golang.org/x/sync/errgroup"
)

// DownlodeClient is the controller of this package
type DownlodeClient struct {
	URL             string
	ContentLength   int64
	HasAcceptRanges bool
	ThreadNumber    int
	IsReady         bool
}

// SetResponceHeader is set response header to DownlodeClient
func (downlodeClient *DownlodeClient) SetResponceHeader() error {
	res, err := http.Head(downlodeClient.URL)
	if err != nil {
		return fmt.Errorf("failed to get Header: %s", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("http error: status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	downlodeClient.ContentLength, err = strconv.ParseInt(res.Header.Get("Content-Length"), 0, 64)
	if err != nil {
		return fmt.Errorf("failed to get content-length: %s", err)
	}

	downlodeClient.HasAcceptRanges = (res.Header.Get("Accept-Ranges") == "bytes")
	downlodeClient.IsReady = downlodeClient.HasAcceptRanges

	return nil
}

func (downlodeClient DownlodeClient) rangeDownload(ctx context.Context, startPos int64, endPos int64) (bytes.Buffer, error) {
	var buf bytes.Buffer
	req, err := http.NewRequest("GET", downlodeClient.URL, nil)
	if err != nil {
		return buf, err
	}
	req = req.WithContext(ctx)

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", startPos, endPos))
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		log.Println("http.Client.Do : ", err)
	}
	defer res.Body.Close()
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		log.Println("rangeDownload error: ", err)
	}
	return buf, nil
}

func writeDownloadData(m map[int][]byte, fileName string) error {

	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		return err
	}
	defer out.Close()

	for i := 0; i <= len(m); i++ {
		_, err = out.Write(m[i])
		if err != nil {
			return err
		}
	}
	return err
}

//Download is download all data by using rangeDownload
func (downlodeClient DownlodeClient) Download(threadNumber int) error {

	if threadNumber == 0 {
		threadNumber = runtime.NumCPU()
	}

	if !downlodeClient.IsReady {
		return fmt.Errorf("DownlodeClient is not ready")
	}

	payloadSize := downlodeClient.ContentLength / int64(threadNumber)
	m := map[int][]byte{}
	remaindSize := downlodeClient.ContentLength

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; remaindSize > 0; i++ {
		i := i
		startPos := downlodeClient.ContentLength - remaindSize
		endPos := startPos + payloadSize - 1
		if endPos > downlodeClient.ContentLength {
			endPos = downlodeClient.ContentLength
		}
		remaindSize -= payloadSize

		eg.Go(func() error {
			buf, err := downlodeClient.rangeDownload(ctx, startPos, endPos)
			m[i] = buf.Bytes()
			return err
		})
	}
	fmt.Println("Download done")
	if err := eg.Wait(); err != nil {
		return err
	}

	return writeDownloadData(m, path.Base(downlodeClient.URL))
}
