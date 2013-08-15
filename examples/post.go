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
	"sync"
)

type NewUrl struct {
	Url string `json:"url"`
}

func readUrls(urls chan string) {
	bio := bufio.NewReaderSize(os.Stdin, 50000)
	for {
		line, isPrefix, err := bio.ReadLine()
		if isPrefix {
			log.Println("uhoh, line too long for buffer ", line)
		} else if err == nil {
			cols := strings.Split(string(line), "\t")
			urls <- cols[2]
		} else if err == io.EOF {
			break
		} else {
			panic(err)
		}
	}
	close(urls)
}

func postUrls(urls chan string) {
	ginger := "http://ec2-54-224-2-249.compute-1.amazonaws.com/collection/wikipedia/"

	for url := range urls {
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

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			postUrls(urls)
		}()
	}
	wg.Wait()
}
