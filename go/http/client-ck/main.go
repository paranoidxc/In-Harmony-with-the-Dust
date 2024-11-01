package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
)

func main() {
	jar, _ := cookiejar.New(nil)

	client := http.Client{
		Jar: jar,
	}

	for i := 0; i < 2; i++ {
		resp, _ := client.Get("http://localhost:18888/cookie")
		dump, _ := httputil.DumpResponse(resp, true)

		log.Println(string(dump))
	}
}
