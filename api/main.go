package main

import (
	"context"
	"net/url"
	"os"
	"os/signal"

	cli "github.com/jawher/mow.cli"
	"github.com/pkg/errors"
	kafka "github.com/segmentio/kafka-go"
)

const appName = "api"

var version string

func main() {
	app := cli.App(appName, "api for all things XCNT")
	app.Version("version", version)

	var (
		debug = app.BoolOpt("debug", false, "print debug logs and enable graphiql")
		port = app.Int(cli.IntOpt{
			Name:   "p port",
			Desc:   "port for the api to listen on",
			EnvVar: "PORT",
			Value:  8080,
		})
		origins = app.Strings(cli.StringsOpt{
			Name:   "origins",
			Desc:   "allowed origins for CORS",
			EnvVar: "ALLOWED_ORIGINS",
			Value:  []string{"http://localhost:8080"},
		})
		dbCxn = app.String(cli.StringOpt{
			Name:   "db",
			Desc:   "DB connection string",
			EnvVar: "DB_CXN",
		})

		kafkaAddr = app.String(cli.StringOpt{
			Name:   "kafka",
			Desc:   "Kafka connection string",
			EnvVar: "KAFKA_ADDR",
		})
		app.Spec = strings.Join([]string{
			"[--debug]",
			"[--port]", "[--origins]",
			"--db",
			"--kafka",
		}, " ")

		app.action = func() {
			logger := logging.NewLogger(appName, *debug)

			store, err := pg.NewStore(*dbCxn)
			if err != nil {
				logger.WithError(err).Fatal("failed to connect to DB")
				cli.Exit(1)
			}

			kafkaURL, err := url.Parse(*kafkaAddr)
			if err != nil {
				logger.WithError(err).Fatal("failed to parse Kafka address")
				cli.Exit(1)
			}

			//TODO create producer

			r := httpsrvr.NewGQLServer(httpsrvr.GQLServerConfig{
				Debug:           *debug,
				Origins:         *origins,
				Logger:          logger,
				TokenSigningKey: *signingKey,
				ExecutableSchema: gql.NewExecutableSchema(gql.Config{
					Resolvers: gql.NewResolver(producers, store),
				}),
			})
		}
	)
}
