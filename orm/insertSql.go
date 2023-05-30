package orm

import (
	"strings"
)

func insertSqlStr(sqlData *SqlData) string {
	sqlname := sqlData.Name
	sqlvalue := sqlData.Value
	index := strings.LastIndex(sqlname, ",")
	if index != -1 {
		sqlname = sqlname[:index]
	}

	index = strings.LastIndex(sqlvalue, ",")
	if index != -1 {
		sqlvalue = sqlvalue[:index]
	}
	return "insert into " + sqlData.Table + " (" + sqlname + ") VALUES (" + sqlvalue + ")"
}

// --- struct to sql
func InsertSqlStr(obj interface{}, params ...OpOption) string {
	op := &Op{sqlType: SQLTYPE_INSERT}
	op.applyOpts(params)
	sqlData := &SqlData{}
	getTableName(obj, sqlData)
	parseStructSql(obj, sqlData, op)
	return insertSqlStr(sqlData)
}

func InsertSql(obj interface{}, params ...OpOption) bool {
	str := InsertSqlStr(obj, params...)
	return exec(str) == nil
}
