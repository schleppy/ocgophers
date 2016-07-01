package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

var (
	errFailedResponseCode = errors.New("400 response")
)

func state(s gobreaker.State) string {
	return []string{"Closed", "Half-Open", "Open"}[s]
}

func main() {
	// START SETUP
	breakerSettings := gobreaker.Settings{
		Name:    "Request local resource",
		Timeout: 5 * time.Second,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("State Change %s --> %s\n", state(from), state(to))
		},
	}
	breakerSettings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		fmt.Printf( // OMIT
			"\tnumReqs %d\n\tfailureRatio %0.2f\n\tconsecutiveFailures %d\n", // OMIT
			counts.Requests,            // OMIT
			failureRatio,               // OMIT
			counts.ConsecutiveFailures, // OMIT
		) // OMIT
		return (counts.Requests > 5 && failureRatio > 0.4) || counts.ConsecutiveFailures > 5
	}
	breakerSettings.MaxRequests = 2
	breaker := gobreaker.NewCircuitBreaker(breakerSettings)
	// END SETUP

	url := "http://localhost:8765"
	for {
		// START CODE
		body, err := breaker.Execute(func() (interface{}, error) {

			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 400 {
				return nil, errFailedResponseCode
			}
			return body, nil
		})
		// END CODE
		time.Sleep(500 * time.Millisecond)
		if err != nil {
			fmt.Printf("Error encountered [%s]\n", err)
			continue
		}
		fmt.Printf("Received body: %s\n", string(body.([]byte)))
	}
}
