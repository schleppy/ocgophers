package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get("timeout")
	to := int64(0)
	var err error
	if timeout != "" {
		to, err = strconv.ParseInt(timeout, 10, 64)
		if err != nil {
			to = 0
		}
	}

	if to > 0 {
		time.Sleep(time.Duration(to) * time.Second)
	}
	if rand.Int()%5 == 0 {
		w.WriteHeader(400)
	}
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":8765", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
