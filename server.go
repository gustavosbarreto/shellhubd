package main

import (
	"context"
	"net"
)

type Server interface {
	CreateListener() (net.Listener, error)
	ConnContext(ctx context.Context, c net.Conn) context.Context
}
