package main

import (
	"context"
	"net/url"
	"os"
	"os/signal"

	cli "github.com/jawher/mow.cli"
	"github.com/pkg/errors"
	kafka "github.com/segmentio/kafka-go"

	"./event"
)

func stripPassword(cxn string) (string, error) {
	u, err := url.Parse(cxn)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse URI")
	}

	if u.User != nil {
		u.User = url.User(u.User.Username())
	}
	return u.String(), nil
}

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

			producers := gql.Producers{
				ExpenseUpdated: event.NewProducer(event.NewKafkaWriter(kafkaURL.Host, appName, event.ExpenseUpdatedTopic)),
				ExpenseInserted: event.NewProducer(event.NewKafkaWriter(
					kafkaURL.Host,
					appName,
					event.ExpenseInsertedTopic)),
			}	

			r := httpsrvr.NewGQLServer(httpsrvr.GQLServerConfig{
				Debug:           *debug,
				Origins:         *origins,
				Logger:          logger,
				TokenSigningKey: *signingKey,
				ExecutableSchema: gql.NewExecutableSchema(gql.Config{
					Resolvers: gql.NewResolver(producers, store),
				}),
			})

			kafkaConn, err := kafka.Dial("tcp", kafkaURL.Host)
			if err != nil {
				logger.WithError(err).Fatal("failed to connect to Kafka")
				cli.Exit(1)
			}

			cleanDBCxn, err := stripPassword(*dbCxn)
			if err != nil {
				logger.WithError(err).Fatal("failed to clean DB cxn string")
				cli.Exit(1)
			}

			cleanKafkaAddr, err := stripPassword(*kafkaAddr)
			if err != nil {
				logger.WithError(err).Fatal("failed to clean Kafka address")
				cli.Exit(1)
			}

			r.Get("/health", health.Handler(health.Settings{
				"version":                  version,
				"debug":                    *debug,
				"appName":                  appName,
				"port":                     *port,
				"databaseConnectionString": cleanDBCxn,
				"kafkaAddress":             cleanKafkaAddr,
			},
				health.PingCheck(logger, "postgres", store, health.Metadata{
					"connectionString": cleanDBCxn,
				}),
				health.Dependency{
					Name: "kafka",
					Metadata: health.Metadata{
						"address": cleanKafkaAddr,
					},
					Check: func(context.Context) bool {
						_, err := kafkaConn.ApiVersions()
						return err == nil
					},
				},
			).ServeHTTP)
			go r.Start(*port)

			stop := make(chan os.Signal)
			signal.Notify(stop, os.Interrupt, os.Kill)
			<-stop
			r.Shutdown()
		}
	
		app.Run(os.Args)
	}
