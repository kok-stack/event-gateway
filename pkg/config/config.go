package config

type EventStoreConfig struct {
	Debug         bool
	Endpoint      string
	SslHost       string
	SslSkipVerify bool
	Verbose       bool
}

type ApplicationConfig struct {
	Test       string
	A          map[string]string
	EventStore EventStoreConfig
}
