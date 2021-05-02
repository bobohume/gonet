 package login

 import (
	 "gonet/base"
	 "gonet/common"
	 "net/http"
 )

type(
	ServerMgr struct{
		m_Inited      bool
		m_config      base.Config
		m_Log         base.CLog
		m_FileMonitor common.IFileMonitor
	}

	IServerMgr interface{
		Init() bool
		GetLog() *base.CLog
		GetFileMonitor() common.IFileMonitor
	}
)

var(
	LoginAddr string
	SERVER ServerMgr
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("login")
	//初始ini配置文件
	this.m_config.Read("GONET_SERVER.CFG")
	LoginAddr = this.m_config.Get3("Login", "Login_Url")

	//动态监控文件改变
	this.m_FileMonitor = &common.FileMonitor{}
	this.m_FileMonitor.Init(1000)

	NETGATECONF.Init()

	http.HandleFunc("/login/", GetNetGateS)
	http.ListenAndServe(LoginAddr, nil)
	return  false
}

 func (this *ServerMgr) GetLog() *base.CLog{
	 return &this.m_Log
 }

 func (this *ServerMgr) GetFileMonitor() common.IFileMonitor {
	 return this.m_FileMonitor
 }