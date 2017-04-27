package main

import (
	"cs/admin"
	"cs/logcs"
	"encoding/json"
	"io/ioutil"
	"ocgophers/domain"
	"ocgophers/profiling/cmd/problem/interfaces"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log"
	"github.com/rcrowley/go-metrics"
)

func main() {

	service := admin.NewAdminService(&admin.AdminServiceConfig{
		AdminPort:   8002,
		ControlPort: 8558,
		AppName:     "WorkHorse",
		GZipCompress: false,
	})

	go metrics.Log(interfaces.MetricsRegistry, 5 * time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	users := loadUsers()

	basicService := interfaces.NewBasicHandler(users)
	counterService := interfaces.NewCounterHandler(users)

	service.Start(basicService, counterService)

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, syscall.SIGTERM, os.Interrupt)
	for range signalChan {
		os.Exit(0)
	}
}

func loadUsers() domain.Responses {
	var responses domain.Responses
	testData, err := ioutil.ReadFile("/opt/crowdstrike/lib/responses.json")
	if err != nil {
		logcs.Logger.Fatalf("Could not open file for reading: [%s]", err)
	}
	err = json.Unmarshal(testData, &responses)
	if err != nil {
		logcs.Logger.Fatalf("Could not unmarshal responses: [%s]", err)
	}

	return responses
}
