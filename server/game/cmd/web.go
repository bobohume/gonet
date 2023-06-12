package cmd

import (
	"gonet/server/cm"
	"net/http"
)

// http://localhost:8080/gm?cmd=cpus()
func cmdHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cmd := r.FormValue("cmd")
	if cmd != "" {
		cm.ParseConsole(g_Cmd, (cmd))
	}
}

func InitWeb() {
	/*go func() {
		http.HandleFunc("/gm", cmdHandle)
		err := http.ListenAndServe(world.Web_Url, nil)
		if err != nil {
			base.LOG.Println("World Web Server : ", err)
		}
	}()*/
}
