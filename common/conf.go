package common

import (
	"gonet/rpc"
)

type (
	Server struct {
		Ip   string `yaml:"ip"`
		Port int    `yaml:"port"`
	}

	Db struct {
		Ip           string `yaml:"ip"`
		Name         string `yaml:"name"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	}

	Redis struct {
		OpenFlag bool   `yaml:"open"`
		Ip       string `yaml:"ip"`
		Password string `yaml:"password"`
	}

	Etcd struct {
		Endpoints []string `yaml:"endpoints"`
	}

	SnowFlake struct {
		Endpoints []string `yaml:"endpoints"`
	}

	Nats struct {
		Endpoints string `yaml:"endpoints"`
	}

	Raft struct {
		Endpoints []string `yaml:"endpoints"`
	}

	Http struct {
		Listen string `yaml:"listen"`
	}

	StubRoute struct {
		STUB rpc.STUB `yaml:"stub"`
	}

	Stub struct {
		StubCount map[string]int64      `yaml:"stub_count"`
		GmCount      int               `yaml:"gm_count"`
	}
)