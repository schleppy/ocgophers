package main

import (
"github.com/CrowdStrike/ratelimiter"
	"fmt"
	"time"
)

func main() {
	maxCapacity := 1000
	ratePeriod := 10 * time.Second
	rl, err := ratelimiter.New(maxCapacity, ratePeriod)
	if err != nil {
		fmt.Printf("Unable to create cache")
	}
	userKey := "sean"
	maxCount := 100 // the maximum number of items I want from this user in ten seconds

	for {
		if cnt, underRateLimit := rl.Incr(userKey, maxCount); underRateLimit {
			fmt.Printf("%s is making request. %d requests made\n", userKey, cnt)
			time.Sleep(50 * time.Millisecond)
		} else {
			fmt.Printf("%s is over rate limit, current count [%d]\n", userKey, cnt)
			time.Sleep(1 * time.Second)
		}
	}

}
