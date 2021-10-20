package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	// External Imports
	"cloud.google.com/go/logging"

	// Internal Imports
	"log-forwarder/types"
)

func HealthCheckHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "I'm healthy this %v at %v:%v:%v", time.Now().UTC().Weekday(), time.Now().UTC().Hour(), time.Now().UTC().Minute(), time.Now().UTC().Second())
}

func PubSubHandler(res http.ResponseWriter, req *http.Request) {
	projectID := fmt.Sprintf("projects/%s", os.Getenv("LOGGING_PROJECT"))

	client, err := logging.NewClient(req.Context(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v \n", err)
	}

	client.OnError = func(err error) {
		log.Printf("client.OnError: %v \n", err)
	}

	defer client.Close()

	name := os.Getenv("LOG_NAME")
	logger := client.Logger(name)

	var m types.PubSubMessage
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v \n", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v \n", err)
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	logMessage := string(m.Message.Data)
	if logMessage == "" {
		log.Printf("Empty Data field in Pub/Sub Message \n")
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}
	logger.Log(logging.Entry{Payload: logMessage})
	res.WriteHeader(http.StatusOK)
}
