package base

import (
	"io/ioutil"
	"log"
	"net"
	"gopkg.in/yaml.v2"
)

func ReadConf(path string, data interface{}) bool{
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("解析config.yaml读取错误: %v", err)
		return  false
	}

	err = yaml.Unmarshal(content, data)
	if err != nil {
		log.Fatalf("解析config.yaml出错: %v", err)
		return false
	}

	return true
}

func GetLanAddr(ip string) (string){
	if ip == "0.0.0.0"{
		addrs, _ := net.InterfaceAddrs()
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					return ip
				}
			}
		}
	}
	return ip
}