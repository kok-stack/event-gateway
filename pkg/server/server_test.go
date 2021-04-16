package server

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"log"
	"testing"
)

func TestSender(t *testing.T) {
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/")

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	e := cloudevents.NewEvent()
	e.SetID("1")
	e.SetType("greeting")
	e.SetSource("test")
	_ = e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
		"message": "Hello, World!",
	})

	res := c.Send(ctx, e)
	if cloudevents.IsUndelivered(res) {
		log.Printf("Failed to send: %v", res)
	} else {
		var httpResult *cehttp.Result
		cloudevents.ResultAs(res, &httpResult)
		log.Printf("%+v", httpResult)
	}
	fmt.Println(res, cloudevents.IsNACK(res))
}
