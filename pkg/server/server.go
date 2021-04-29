package server

import (
	"context"
	tracingclient "github.com/cloudevents/sdk-go/observability/opencensus/v2/client"
	obshttp "github.com/cloudevents/sdk-go/observability/opencensus/v2/http"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol"
	client3 "github.com/jdextraze/go-gesclient/client"
	"github.com/kok-stack/event-gateway/pkg/config"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
)

var dbClient client3.Connection

func StartServer(ctx context.Context, config *config.ApplicationConfig, dc client3.Connection) error {
	dbClient = dc
	p, err := obshttp.NewObservedHTTP()
	if err != nil {
		return err
	}
	c, err := client.New(p, client.WithObservabilityService(tracingclient.New()))
	if err != nil {
		return err
	}
	go func() {
		select {
		case <-ctx.Done():
			dbClient.Close()
		}
	}()

	return c.StartReceiver(ctx, handler)
}

func handler(ctx context.Context, e event.Event) protocol.Result {
	log.Printf("-> \n %+v", e.String())

	v4, _ := uuid.NewV4()
	evt := client3.NewEventData(v4, e.Type(), true, e.Data(), nil)
	var err error
	val, err := strconv.Atoi(e.ID())
	if err != nil {
		return protocol.NewReceipt(false, "convert id to version error: %v", err)
	}
	if val <= 0 {
		val = client3.ExpectedVersion_Any
	}

	log.Printf("-> '%s': %+v", "test", evt)

	task, err := dbClient.AppendToStreamAsync(e.Subject(), val, []*client3.EventData{evt}, nil)
	if err != nil {
		log.Printf("Error occured while appending to stream: %v", err)
		return protocol.NewReceipt(false, "Error occured while appending to stream: %v", err)
	} else if err := task.Error(); err != nil {
		log.Printf("Error occured while waiting for result of appending to stream: %v", err)
		return protocol.NewReceipt(false, "Error occured while waiting for result of appending to stream: %v", err)
	} else {
		result := task.Result().(*client3.WriteResult)
		log.Printf("<- %+v", result)
	}
	return protocol.ResultACK
}
