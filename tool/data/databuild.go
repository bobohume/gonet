package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	args := os.Args
	waitGroup := &sync.WaitGroup{}
	//获取当前路径
	//dir, _ := filepath.Abs(`.`)
	if len(args) >= 2 && args[1] == "decode" {
		files, err := filepath.Glob("*.dat")
		if err == nil {
			for _, v := range files {
				waitGroup.Add(1)
				go func(name string) {
					SaveExcel(name)
					waitGroup.Done()
				}(v)
			}
		}
	} else {
		files, err := filepath.Glob("*.xlsx")
		if err == nil {
			for _, v := range files {
				waitGroup.Add(1)
				go func(name string) {
					//OpenExcel(name)
					OpenExceLua(name)
					OpenExceGo(name)
					//OpenExceCsv(name)
					waitGroup.Done()
				}(v)
			}
		}
	}
	waitGroup.Wait()
	fmt.Println("解析完成")
}
