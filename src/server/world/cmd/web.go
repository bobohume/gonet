package cmd

import (
	"net/http"
	"strings"
	"server/common"
	"server/world"
)

func webHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//http://localhost:8080/gm?cmd=cpus()
	if r.URL.Path == "/gm"{
		cmd := r.Form["cmd"]
		call := strings.Join(cmd, "")
		call = strings.TrimSpace(call)
		common.ParseConsole(g_Cmd, []byte(call))
	}
}

func InitWeb(){
	go func() {
		http.HandleFunc("/", webHandle)
		err := http.ListenAndServe(world.Web_Url, nil)
		if err != nil {
			world.SERVER.GetLog().Println("World Web Server : ", err)
		}
	}()
}
