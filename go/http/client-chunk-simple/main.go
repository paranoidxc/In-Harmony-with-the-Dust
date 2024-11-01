package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
)

func main() {
	resp, _ := http.Get("http://localhost:18888/chunked")
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		log.Println(string(line))
	}
}
