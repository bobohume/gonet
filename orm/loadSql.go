package orm

import (
	"gonet/base"
	"strings"
)

func loadSqlStr(sqlData *SqlData, op *Op) string {
	sqlname := sqlData.Name
	key := sqlData.Key
	index := strings.LastIndex(sqlname, ",")
	if index != -1 {
		sqlname = sqlname[:index]
	}
	if !op.whereFlag {
		index = strings.LastIndex(key, ",")
		if index != -1 {
			key = " where " + key[:index]
		}
	} else if op.where != "" {
		key = " where " + op.where
	} else {
		key = ""
	}
	return "select " + sqlname + " from " + sqlData.Table + key + op.limit
}

// --- struct to sql
func LoadSqlStr(obj interface{}, params ...OpOption) string {
	op := &Op{sqlType: SQLTYPE_LOAD}
	op.applyOpts(params)
	sqlData := &SqlData{}
	getTableName(obj, sqlData)
	parseStructSql(obj, sqlData, op)
	return loadSqlStr(sqlData, op)
}

func LoadSql(obj interface{}, params ...OpOption) (*Rows, error) {
	str := LoadSqlStr(obj, params...)
	rs, err := Query(str)
	if err != nil {
		base.LOG.Fatalf("%s %s", str, err.Error())
	}
	return rs, err
}
