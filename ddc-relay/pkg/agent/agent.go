package agent

import (
	"os"
	"os/signal"
	"time"

	"github.com/NativeMetrics/ddc-relayy/pkg/settings"
	"github.com/nats-io/nats.go"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

func Run() error {
	// connect to nats
	nc, err := Connect()
	if err != nil {
		return eris.Wrap(err, "failed to connect to nats")
	}
	defer nc.Close()

	// create subscriptions
	if err = createSubscriptions(nc); err != nil {
		return eris.Wrap(err, "failed to create subscriptions")
	}

	// flush NATS connection
	if err = nc.Flush(); err != nil {
		return eris.Wrap(err, "failed to flush nats connection")
	}

	// check last error
	if err := nc.LastError(); err != nil {
		return eris.Wrap(err, "failed to get last error")
	}

	// listen for messages
	log.Info().Msg("connected to NATS, listening for messages")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Info().Msg("draining nats connection")

	if err = nc.Drain(); err != nil {
		return eris.Wrap(err, "failed to drain nats connection")
	}

	return nil
}

const (
	totalWait      = 10 * time.Minute
	reconnectDelay = time.Second
)

func Connect() (*nats.Conn, error) {
	opts := []nats.Option{
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(int(totalWait / reconnectDelay)),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Error().
				Err(err).
				Float64("reconnectWindow", totalWait.Minutes()).
				Msg("nats disconnected")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info().
				Str("url", nc.ConnectedUrl()).
				Msg("nats reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Warn().
				Err(nc.LastError()).
				Msg("nats exiting")
		}),
	}

	n, err := nats.Connect(settings.Values.NatsURL, opts...)
	if err != nil {
		return nil, eris.Wrap(err, "failed trying to connect to nats server")
	}

	return n, nil
}
