package common

import (
	"base"
)

/*const(
	NONE_ERROR		=iota,
)*/

func DBERROR(msg string, err error){
	base.GLOG.Printf("db [%s] error [%s]", msg, err.Error())
}
