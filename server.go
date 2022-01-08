package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"example.com/goserver/app"
)

var counter int
var mutex = &sync.Mutex{}

func appMeta(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		reqBody, err := ioutil.ReadAll(req.Body)

		if err != nil {
			log.Fatal(err)
		}

		var meta app.Meta
		err = meta.Parse(reqBody)
		if err != nil {
			log.Fatal("Failed to parse request yaml")
		}
		fmt.Fprintf(w, "%s\n%s", meta.Title, meta.Version)
	}
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
}

func queryParams(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func main() {
	http.HandleFunc("/apps", appMeta)
	http.HandleFunc("/increment", incrementCounter)
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
