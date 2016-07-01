package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eapache/go-resiliency/retrier"
)

var (
	errFailedResponseCode = errors.New("400 response")
)

func main() {
	// START SETUP
	retry := retrier.New(retrier.ConstantBackoff(2, 10*time.Millisecond), nil)
	// END SETUP

	url := "http://localhost:8765"
	tempFailures := 0
	failures := 0
	success := 0
	for {
		// START CODE
		var body []byte
		reqErr := retry.Run(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, _ = ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 400 {
				tempFailures++
				return errFailedResponseCode
			}
			return nil
		})
		// END CODE
		time.Sleep(250 * time.Millisecond)
		if reqErr != nil {
			failures++
			fmt.Printf("Error encountered [%s]\n", reqErr)
			continue
		}
		success++
		fmt.Printf("Temporary failures: %d, Hard failures: %d, Successes: %d\n", tempFailures, failures, success)
	}
}
