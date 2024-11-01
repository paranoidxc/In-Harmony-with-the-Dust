package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func readFromString() *strings.Reader {
	reader := strings.NewReader(" heesdsdfsfdsf")

	return reader
}

func main() {
	//file, _ := os.Open("./main.go")
	file := readFromString()
	resp, err := http.Post("http://localhost:18888", "text/plain", file)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Println("header", resp.Header)
	log.Println("Content-Length ", resp.Header.Get("Content-Length"))

	log.Println("Status", resp.Status)
	log.Println("StatusCode", resp.StatusCode)
	log.Println(string(body))
}
