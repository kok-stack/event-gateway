package db

import (
	"github.com/jdextraze/go-gesclient"
	client2 "github.com/jdextraze/go-gesclient/client"
	"github.com/kok-stack/event-gateway/pkg/config"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

func ConnTCP(config *config.ApplicationConfig) (client2.Connection, error) {
	//gesclient.Debug()

	c, err := createConnection(config)
	if err != nil {
		panic(err.Error())
	}
	c.Connected().Add(func(evt client2.Event) error { log.Printf("Connected: %+v", evt); return nil })
	c.Disconnected().Add(func(evt client2.Event) error { log.Printf("Disconnected: %+v", evt); return nil })
	c.Reconnecting().Add(func(evt client2.Event) error { log.Printf("Reconnecting: %+v", evt); return nil })
	c.Closed().Add(func(evt client2.Event) error { log.Fatalf("Connection closed: %+v", evt); return nil })
	c.ErrorOccurred().Add(func(evt client2.Event) error { log.Printf("Error: %+v", evt); return nil })
	c.AuthenticationFailed().Add(func(evt client2.Event) error { log.Printf("Auth failed: %+v", evt); return nil })

	if err := c.ConnectAsync().Wait(); err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	return c, nil
}

func createConnection(c *config.ApplicationConfig) (client2.Connection, error) {
	settingsBuilder := client2.CreateConnectionSettings()

	var uri *url.URL
	var err error
	if !strings.Contains(c.EventStore.Endpoint, "://") {
		gossipSeeds := strings.Split(c.EventStore.Endpoint, ",")
		endpoints := make([]*net.TCPAddr, len(gossipSeeds))
		for i, gossipSeed := range gossipSeeds {
			endpoints[i], err = net.ResolveTCPAddr("tcp", gossipSeed)
			if err != nil {
				log.Fatalf("Error resolving: %v", gossipSeed)
			}
		}
		settingsBuilder.SetGossipSeedEndPoints(endpoints)
	} else {
		uri, err = url.Parse(c.EventStore.Endpoint)
		if err != nil {
			log.Fatalf("Error parsing address: %v", err)
		}

		if uri.User != nil {
			username := uri.User.Username()
			password, _ := uri.User.Password()
			settingsBuilder.SetDefaultUserCredentials(client2.NewUserCredentials(username, password))
		}
	}

	if c.EventStore.SslHost != "" {
		settingsBuilder.UseSslConnection(c.EventStore.SslHost, !c.EventStore.SslSkipVerify)
	}

	if c.EventStore.Verbose {
		settingsBuilder.EnableVerboseLogging()
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("get hostname Error: %v", err)
	}

	return gesclient.Create(settingsBuilder.Build(), uri, hostname)
}
