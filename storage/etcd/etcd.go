package etcd

import (
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	Endpoints   []string
	DialTimeout time.Duration
}

func DefaultConfigs() *Config {
	return &Config{
		DialTimeout: 5 * time.Second,
	}
}

type etcd struct {
	client *clientv3.Client
	log    *log.Logger
}

func New(log *log.Logger, cfg Config) (*etcd, error) {
	return newClient(log, &cfg)
}

func newClient(log *log.Logger, cfg *Config) (*etcd, error) {
	if cfg == nil {
		cfg = DefaultConfigs()
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
	})

	if err != nil {
		return nil, err
	}

	return &etcd{client: cli, log: log}, nil
}

func (etcd etcd) Close() {
	etcd.client.Close()
}
