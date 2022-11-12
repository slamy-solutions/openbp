package nats

import (
	"github.com/nats-io/nats.go"
)

func Connect(url string, name string) (*nats.Conn, error) {
	return nats.Connect(url, nats.Name(name))
}
