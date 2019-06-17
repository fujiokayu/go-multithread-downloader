package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// MyHandler is a object of this http server
type MyHandler struct {
}

func (MyHandler *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]
	log.Println(string(path))
	fmt.Println(string(path))
	if string(path) == "" {
		http.StripPrefix("/", http.FileServer(http.Dir("./testdata/contents")))
	} else {
		log.Println("download file")
		data, err := ioutil.ReadFile(string(path))
		if err != nil {
			log.Println(err)
		} else {
			info, err := os.Stat(string(path))
			if err != nil {
				log.Println("there was an error to get file stats.", err)
			}

			w.Header().Set("Content-Disposition", "attachment; filename="+string(path))
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
			w.Write(data)
		}
	}
}

func main() {
	http.Handle("/", new(MyHandler))

	fmt.Println("URL: http://localhost:8080/contents/")
	http.ListenAndServe(":8080", nil)
}
