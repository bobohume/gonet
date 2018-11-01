package main

import (
	"log"
	"net/http"
	"strings"
	"os/exec"
	"fmt"
	"bytes"
	"io"
)

//for world
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//http://localhost:8080/gm?cmd=cpus()
	if r.URL.Path == "/gm"{
		cmd := r.Form["cmd"]
		call := strings.Join(cmd, "")
		call = strings.TrimSpace(call)
		if call == "restart" {
			cmd := exec.Command("sh", "-c", "/root/restart.sh")
			cmd.Stdin = strings.NewReader("some input")
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Run()
			fmt.Printf(out.String())
			io.WriteString(w, out.String())
			io.WriteString(w, "-----------重启完毕--------------")
		}else if  call == "update"{
			{
				cmd := exec.Command("ssh", "root@192.168.5.10", "cd /www/newmir/res/ && svn up")
				cmd.Stdin = strings.NewReader("some input")
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()
				fmt.Println(out.String())
				io.WriteString(w, out.String())
			}
			{
				cmd := exec.Command("ssh", "root@192.168.5.10", "cd /root/swy/gameserver/bin/ && svn up")
				cmd.Stdin = strings.NewReader("some input")
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()
				fmt.Println(out.String())
				io.WriteString(w, out.String())
			}
			{
				cmd := exec.Command("sh", "-c", "/root/go/update_res.sh")
				cmd.Stdin = strings.NewReader("some input")
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()
				fmt.Println(out.String())
				io.WriteString(w, out.String())
				io.WriteString(w, "-----------更新资源完毕--------------")
			}
		} else if  call == "update_code"{
			{
				cmd := exec.Command("ssh", "root@192.168.5.10", "cd /root/swy/gameserver && svn up")
				cmd.Stdin = strings.NewReader("some input")
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()
				fmt.Println(out.String())
				io.WriteString(w, out.String())
			}

			{
				cmd := exec.Command("ssh", "root@192.168.5.10", "cd /root/swy/gameserver/scripts &&./build.sh")
				cmd.Stdin = strings.NewReader("some input")
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()
				fmt.Println(out.String())
				io.WriteString(w, out.String())
			}

			{
				cmd := exec.Command("ssh", "root@192.168.5.10", "cd /root/swy && scp_to_511.sh")
				cmd.Stdin = strings.NewReader("some input")
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()
				fmt.Println(out.String())
				io.WriteString(w, out.String())
				io.WriteString(w, "-----------编译代码完毕--------------")
			}
		}else if call == "date" {
			cmd1 := r.Form["time"]
			time := strings.Join(cmd1, "")
			time = strings.TrimSpace(time)

			cmd := exec.Command("date", "-s", time)
			cmd.Stdin = strings.NewReader("some input")
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Run()
			fmt.Printf(out.String())
			io.WriteString(w, out.String())
			io.WriteString(w, "-----------时间设置成功--------------")
		}
	}
}
func main() {
	http.HandleFunc("/", sayhelloName)
	err := http.ListenAndServe("192.168.5.11:9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}