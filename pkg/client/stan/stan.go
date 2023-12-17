package stan

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type client struct {
	url       string
	clusterID string
}

// New returns a nats streaming driver
func New(url, clusterID string) *client {
	return &client{
		url:       url,
		clusterID: clusterID,
	}
}

func (c *client) NewConnection(clientID string) (stan.Conn, error) {
	opts := []nats.Option{
		nats.NoReconnect(),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("nats connection lost: %s\n", err.Error())
		}),
		nats.ReconnectHandler(func(nnc *nats.Conn) {
			fmt.Println("reconnected to NATS server")
		}),
	}
	nc, err := nats.Connect(c.url, opts...)
	if err != nil {
		fmt.Printf("failed to connect to NATS server: %s\n", err.Error())
		return nil, err
	}
	sc, err := stan.Connect(c.clusterID, clientID, stan.NatsConn(nc), stan.Pings(5, 60),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			fmt.Printf("NATS streaming connection lost: %s", err.Error())
		}))
	if err != nil {
		fmt.Printf("failed to connect to NATS streaming server: %s", err.Error())
		return nil, err
	}
	return sc, nil
}
