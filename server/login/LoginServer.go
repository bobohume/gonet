package login

import (
	"gonet/base"
	"gonet/common"
	"net/http"
)

type (
	ServerMgr struct {
		isInited    bool
		fileMonitor common.IFileMonitor
	}

	IServerMgr interface {
		Init() bool
		GetFileMonitor() common.IFileMonitor
	}

	Config struct {
		common.Http `yaml:"login"`
	}
)

var (
	CONF   Config
	SERVER ServerMgr
)

func (s *ServerMgr) Init() bool {
	if s.isInited {
		return true
	}

	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	//动态监控文件改变
	s.fileMonitor = &common.FileMonitor{}
	s.fileMonitor.Init()

	NETGATECONF.Init()

	http.HandleFunc("/login/", GetNetGateS)
	http.ListenAndServe(CONF.Http.Listen, nil)
	return false
}

func (s *ServerMgr) GetFileMonitor() common.IFileMonitor {
	return s.fileMonitor
}
