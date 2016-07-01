package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/net/context"
)

var (
	wg = sync.WaitGroup{}
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	result := make(chan int, 2)
	wg.Add(1)
	go doSomething(ctx, result)
	select {
	case <-ctx.Done():
		fmt.Println("We give up")
	case c := <-result:
		fmt.Println("Work complete.  Answer is", c)
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
}

func doSomething(ctx context.Context, result chan int) {
	defer wg.Done()
	t := 200
	fmt.Println("time to wait", t)
	select {
	case <-time.After(time.Duration(t) * time.Millisecond):
		fmt.Println("doing something")
		result <- 42
	case <-ctx.Done():
		fmt.Println("We already gave up")
	}
}
