package orm

import (
	"strings"
)

func saveSqlStr(sqlData *SqlData) string {
	sqlname := sqlData.Name
	sqlvalue := sqlData.Value
	sqlext := sqlData.NameValue
	index := strings.LastIndex(sqlname, ",")
	if index != -1 {
		sqlname = sqlname[:index]
	}

	index = strings.LastIndex(sqlvalue, ",")
	if index != -1 {
		sqlvalue = sqlvalue[:index]
	}
	index = strings.LastIndex(sqlext, ",")
	if index != -1 {
		sqlext = sqlext[:index]
	}
	return "insert into " + sqlData.Table + " (" + sqlname + ") VALUES (" + sqlvalue + ") ON DUPLICATE KEY UPDATE" + sqlext
}

// --- struct to sql
func SaveSqlStr(obj interface{}, params ...OpOption) string {
	op := &Op{sqlType: SQLTYPE_SAVE}
	op.applyOpts(params)
	sqlData := &SqlData{}
	getTableName(obj, sqlData)
	parseStructSql(obj, sqlData, op)
	return saveSqlStr(sqlData)
}

func SaveSql(obj interface{}, params ...OpOption) bool {
	str := SaveSqlStr(obj, params...)
	return exec(str) == nil
}
