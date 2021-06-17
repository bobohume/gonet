package login

import (
	"gonet/base"
	"gonet/base/ini"
	"net/http"
	"sync"
)

type(
	NetGateConf struct{
		m_config ini.Config
		m_Locker *sync.RWMutex
	}
)

var(
	NETGATECONF	NetGateConf
)

func (this *NetGateConf) Init() bool {
	this.m_Locker = &sync.RWMutex{}
	this.Read()
	SERVER.GetFileMonitor().AddFile("NETGATES.CFG", this.Read)
	return true
}

func (this *NetGateConf) Read() {
	this.m_Locker.Lock()
	this.m_config.Read("NETGATES.CFG")
	this.m_Locker.Unlock()
}

func (this *NetGateConf) GetNetGates(Arena string) []string{
	this.m_Locker.RLock()
	arenas := this.m_config.Get6(Arena, "NetGates", ",")
	this.m_Locker.RUnlock()
	return arenas
}

func GetNetGateS(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	arenas := NETGATECONF.GetNetGates(r.FormValue("arena"))
	nLen := len(arenas)
	if nLen > 0{
		nIndex := base.RAND.RandI(0, nLen-1)
		w.Write([]byte(arenas[nIndex]))
		return
	}

	w.Write([]byte(""))
}