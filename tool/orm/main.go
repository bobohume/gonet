package main

import "gonet/server/model"

//args[1] : proto file path
func main(){
	model.Generate("Player")
}
