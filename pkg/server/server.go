package server

import (
	"context"
	"fmt"
	tracingclient "github.com/cloudevents/sdk-go/observability/opencensus/v2/client"
	obshttp "github.com/cloudevents/sdk-go/observability/opencensus/v2/http"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/go-redis/redis/v8"
	"github.com/kok-stack/event-gateway/pkg/config"
)

var dbClient *redis.Client

func StartServer(ctx context.Context, config *config.ApplicationConfig, dc *redis.Client) error {
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
	fmt.Println(e.String())
	add := dbClient.PFAdd(ctx, e.Type(), e.ID())
	result, err := add.Result()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if result <= 0 {
		return protocol.NewReceipt(false, "event %s id:%s exists", e.Type(), e.ID())
	}

	return nil
}
