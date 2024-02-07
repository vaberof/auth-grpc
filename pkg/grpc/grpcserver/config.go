package grpcserver

import "time"

type ServerConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
}
