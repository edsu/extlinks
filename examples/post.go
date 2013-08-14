package main

// post.go is a golang program that reads the output of parse.py
// and POSTs the URLs as JSON to a web server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type NewUrl struct {
	Url string `json:"url"`
}

func readUrls(urls chan string) {
	bio := bufio.NewReaderSize(os.Stdin, 50000)
	for {
		line, isPrefix, err := bio.ReadLine()
		if isPrefix {
			// url was too big. m'eh who cares for now ...
		} else if err == nil {
			cols := strings.Split(string(line), "\t")
			urls <- cols[2]
		} else if err == io.EOF {
			break
		} else {
			panic(err)
		}
	}
}

func postUrls(urls chan string) {
	ginger := "http://example.com/collection/wikipedia/"

	for {
		url := <-urls
		n := NewUrl{url}
		data, _ := json.Marshal(n)
		resp, err := http.Post(ginger, "application/json", bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		} else if resp.StatusCode != http.StatusOK {
			log.Fatal("received ", resp.Status)
		} else {
			log.Println("added ", url)
		}
	}
}

func main() {
	urls := make(chan string)
	go readUrls(urls)
	postUrls(urls)
}
