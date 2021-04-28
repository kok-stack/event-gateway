package main

import (
	"context"
	"fmt"
	config2 "github.com/kok-stack/event-gateway/pkg/config"
	"github.com/kok-stack/event-gateway/pkg/db"
	"github.com/kok-stack/event-gateway/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func NewCommand() (*cobra.Command, context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	config := &config2.ApplicationConfig{}
	command := &cobra.Command{
		Use:   ``,
		Short: ``,
		Long:  ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Kill)
				<-c
				cancelFunc()
			}()

			err := viper.Unmarshal(config)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%+v \n", config)
			config.EventStore = config2.EventStoreConfig{
				Debug:         true,
				Endpoint:      "tcp://admin:changeit@127.0.0.1:1113",
				SslHost:       "",
				SslSkipVerify: false,
				Verbose:       false,
			}

			conn, err := db.ConnTCP(config)
			if err != nil {
				return err
			}

			//启动cloudevents
			err = server.StartServer(ctx, config, conn)
			if err != nil {
				return err
			}

			<-ctx.Done()
			return nil
		},
	}
	//command.AddCommand()
	viper.AutomaticEnv()
	viper.AddConfigPath(`.`)
	command.PersistentFlags().String("test", "abcd", "abcd")
	command.PersistentFlags().StringToString("a", nil, "")
	err := viper.BindPFlags(command.PersistentFlags())
	if err != nil {
		panic(err.Error())
	}
	err = viper.BindPFlags(command.Flags())
	if err != nil {
		panic(err.Error())
	}
	viper.SetDefault("test", "dcba")

	return command, ctx, cancelFunc
}

func main() {
	command, _, _ := NewCommand()
	if err := command.Execute(); err != nil {
		panic(err)
	}
}
