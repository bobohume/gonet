package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func main(){
	args := os.Args
	waitGroup := &sync.WaitGroup{}
	if len(args) >= 2 && args[1] == "decode"{
		files1, err := filepath.Glob("*.dat")
		if err == nil{
			for _, v := range files1{
				waitGroup.Add(1)
				go func(name string) {
					SaveExcel(name)
					waitGroup.Done()
				}(v)
			}
		}
	}else{
		files, err := filepath.Glob("*.xlsx")
		if err == nil{
			for _, v := range files{
				waitGroup.Add(1)
				go func(name string) {
					OpenExcel(name)
					waitGroup.Done()
				}(v)
			}
		}
	}
	waitGroup.Wait()
	fmt.Println("解析完成")
}