package login

import (
	"gonet/base"
	"gonet/base/ini"
	"net/http"
	"sync"
)

type (
	NetGateConf struct {
		config ini.Config
		locker *sync.RWMutex
	}
)

var (
	NETGATECONF NetGateConf
)

func (n *NetGateConf) Init() bool {
	n.locker = &sync.RWMutex{}
	n.Read()
	SERVER.GetFileMonitor().AddFile("NETGATES.CFG", n.Read)
	return true
}

func (n *NetGateConf) Read() {
	n.locker.Lock()
	n.config.Read("NETGATES.CFG")
	n.locker.Unlock()
}

func (n *NetGateConf) GetNetGates(Arena string) []string {
	n.locker.RLock()
	arenas := n.config.Get6(Arena, "NetGates", ",")
	n.locker.RUnlock()
	return arenas
}

func GetNetGateS(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	arenas := NETGATECONF.GetNetGates(r.FormValue("arena"))
	nLen := len(arenas)
	if nLen > 0 {
		nIndex := base.RandI(0, nLen-1)
		w.Write([]byte(arenas[nIndex]))
		return
	}

	w.Write([]byte(""))
}
