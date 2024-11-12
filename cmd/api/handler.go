package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/obynonwane/broker-service/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type RPCPayload struct {
	Name string
	Data string
}

type RequestPayload struct {
	Action string `json:"action"`
}

func (app *Config) TestSetup(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome")
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// leave connection open forever
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// leave cancel
	defer cancel()

	data, err := c.WriteLog(ctx, &logs.LogRequest{LogEntry: &logs.Log{
		Name: "Testing out",
		Data: "Testing data",
	}})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	log.Println("the logged paylod", payload)
	log.Println("the logged paylod the data", data)

	app.writeJSON(w, http.StatusAccepted, payload)

}
