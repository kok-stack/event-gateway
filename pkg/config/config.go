package config

type EventStoreConfig struct {
	Debug         bool
	Endpoint      string
	SslHost       string
	SslSkipVerify bool
	Verbose       bool
}

type BridgeConfig struct {
	StreamName string
	Group      string
}

type ApplicationConfig struct {
	Test       string
	A          map[string]string
	EventStore EventStoreConfig
	Bridge     BridgeConfig
}
