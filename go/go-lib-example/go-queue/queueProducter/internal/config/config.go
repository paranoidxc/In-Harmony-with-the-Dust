package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
