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
		StubStrRoute map[string]string `yaml:"stub_route"`
		StubStrCount map[string]int    `yaml:"stub_count"`
		GmCount      int               `yaml:"gm_count"`
		StubRoute    map[string]rpc.STUB
		StubCount    map[rpc.STUB]int
	}
)

func (this *Stub) Init() {
	this.StubRoute = map[string]rpc.STUB{}
	this.StubCount = map[rpc.STUB]int{}
	for k, v := range this.StubStrRoute {
		this.StubRoute[k] = rpc.STUB(rpc.STUB_value[v])
	}

	for k, v := range this.StubStrCount {
		this.StubCount[rpc.STUB(rpc.STUB_value[k])] = v
	}
}
