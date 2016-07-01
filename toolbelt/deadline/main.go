package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/eapache/go-resiliency/deadline"
)

var (
	errFailedResponseCode = errors.New("400 response")
)

func main() {
	// START SETUP
	dl := deadline.New(1 * time.Second)
	// END SETUP

	for {
		var body []byte
		url := "http://localhost:8765"
		if rand.Int()%5 == 0 {
			url = url + "?timeout=2"
		}
		tStart := time.Now()
		// START CODE
		err := dl.Run(func(stopper <-chan struct{}) error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, _ = ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 400 {
				return errFailedResponseCode
			}
			return nil
		})

		delta := time.Since(tStart).Nanoseconds() / 1e6
		switch err {
		case deadline.ErrTimedOut:
			fmt.Printf("Timeout error: %d ms\n", delta)
		case nil:
			fmt.Printf("Request response: %s, %d ms\n", string(body), delta)
		default:
			fmt.Printf("Some other error: %s, %d ms\n", err, delta)
		}
		// END CODE
		time.Sleep(500 * time.Millisecond)
	}
}
