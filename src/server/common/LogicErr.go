package common

import (
	"server/world"
)

/*const(
	NONE_ERROR		=iota,
)*/

func DBERROR(msg string, err error){
	world.SERVER.GetLog().Printf("db [%s] error [%s]", msg, err.Error())
}
