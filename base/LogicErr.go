package base

/*const(
	NONE_ERROR		=iota,
)*/

func DBERROR(msg string, err error) {
	LOG.Printf("db [%s] error [%s]", msg, err.Error())
}
