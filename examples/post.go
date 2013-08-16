package main

// post.go is a golang program that reads the output of parse.py
// and POSTs the URLs as JSON to a web server

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"code.google.com/p/go.net/websocket"
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
	close(urls)
}

func postUrls(urls chan string) {
	for url := range urls {
		n := NewUrl{url}
		data, _ := json.Marshal(n)
		resp, err := http.Post(*ginger, "application/json", bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		} else if resp.StatusCode == http.StatusCreated {
			log.Println("added ", url)
		} else {
			log.Println("received ", resp.Status)
		}
	}
}

func sendUrls(urls chan string) {
	origin := "http://" + *ginger + "/"
	url := "ws://" + *ginger + ":80/collection/one/add"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	for url := range urls {
		websocket.Message.Send(ws, url)

		var message string
		websocket.Message.Receive(ws, &message)
		log.Println(message)
	}
}

var ginger = flag.String("ginger", "localhost", "url to ginger collection")

func main() {
	flag.Parse()
	urls := make(chan string)
	go readUrls(urls)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendUrls(urls)
		}()
	}
	wg.Wait()
}
