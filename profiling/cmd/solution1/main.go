package main

import (
	"cs/admin"
	"cs/logcs"
	"encoding/json"
	"io/ioutil"
	"ocgophers/domain"
	"ocgophers/profiling/cmd/solution1/interfaces"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	service := admin.NewAdminService(&admin.AdminServiceConfig{
		AdminPort:   8001,
		ControlPort: 8558,
		AppName:     "WorkHorse",
		GZipCompress: false,
	})

	users := loadUsers()

	basicService := interfaces.NewBasicHandler(users)

	service.Start(basicService)

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, syscall.SIGTERM, os.Interrupt)
	for range signalChan {
		os.Exit(0)
	}
}

func loadUsers() [][]byte {
	var responses domain.Responses
	testData, err := ioutil.ReadFile("/opt/crowdstrike/lib/responses.json")
	if err != nil {
		logcs.Logger.Fatalf("Could not open file for reading: [%s]", err)
	}
	err = json.Unmarshal(testData, &responses)
	if err != nil {
		logcs.Logger.Fatalf("Could not unmarshal responses: [%s]", err)
	}

	readyResponses := [][]byte{}
	for i, response := range responses {
		d, err := json.Marshal(response)
		if err != nil {
			logcs.Logger.Errorf("Could not marshal record %d", i)
			continue
		}
		readyResponses = append(readyResponses, d)
	}
	return readyResponses
}
