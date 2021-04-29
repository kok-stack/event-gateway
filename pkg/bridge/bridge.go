package bridge

import (
	"context"
	"github.com/jdextraze/go-gesclient/client"
	"github.com/kok-stack/event-gateway/pkg/config"
	"log"
)

func StartBridge(ctx context.Context, cli client.Connection, cfg *config.ApplicationConfig) error {
	task, err := cli.ConnectToPersistentSubscriptionAsync(cfg.Bridge.StreamName, cfg.Bridge.Group, appendHandler, dropHandler, nil, 10, false)
	if err != nil {
		//create
		task, err = cli.CreatePersistentSubscriptionAsync(cfg.Bridge.StreamName, cfg.Bridge.Group, client.DefaultPersistentSubscriptionSettings, nil)
		if err != nil {
			log.Printf("Error occured while subscribing to stream: %v", err)
			return err
		} else if err := task.Error(); err != nil {
			log.Printf("Error occured while waiting for result of subscribing to stream: %v", err)
			return err
		}
		res := task.Result().(*client.PersistentSubscriptionCreateResult)
		if res.GetStatus() != client.PersistentSubscriptionCreateStatus_Success {
			log.Printf("CreatePersistentSubscriptionAsync faild: %+v", res)
		}
		log.Printf("Error occured while ConnectToPersistentSubscriptionAsync to stream: %v", err)
	} else if err := task.Error(); err != nil {
		log.Printf("Error occured while waiting for result of subscribing to stream: %v", err)
	} else {
		sub := task.Result().(client.PersistentSubscription)
		log.Printf("SubscribeToStream result: %+v", sub)
		go func() {
			select {
			case <-ctx.Done():
				sub.Stop()
			}
		}()
	}

	return nil
}

func dropHandler(s client.PersistentSubscription, dr client.SubscriptionDropReason, err error) error {
	return nil
}

func appendHandler(s client.PersistentSubscription, r *client.ResolvedEvent) error {
	//转发event到url(cloud event协议)
	//r.Event().
	//s.Acknowledge()
	return nil
}
