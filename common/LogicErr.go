package common

import (
	"gonet/base"
)

/*const(
	NONE_ERROR		=iota,
)*/

func DBERROR(msg string, err error) {
	base.LOG.Printf("db [%s] error [%s]", msg, err.Error())
}
