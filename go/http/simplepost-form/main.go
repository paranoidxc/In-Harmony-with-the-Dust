package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("name", "Michael j")
	fileWriter, err := writer.CreateFormFile("thumbnial", "./main.go")
	if err != nil {
		panic(err)
	}
	readFile, _ := os.Open("./main.go")
	defer readFile.Close()

	io.Copy(fileWriter, readFile)
	writer.Close()

	resp, err := http.Post("http://localhost:18888", writer.FormDataContentType(), &buffer)
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
