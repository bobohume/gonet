package orm

import (
	"strings"
)

func deleteSqlStr(sqlData *SqlData) string {
	key := sqlData.Key
	index := strings.LastIndex(key, ",")
	if index != -1 {
		key = key[:index]
	}
	key = strings.Replace(key, ",", " and ", -1)
	return "delete from " + sqlData.Table + " where " + key
}

// --- struct to sql
func DeleteSqlStr(obj interface{}, params ...OpOption) string {
	op := &Op{sqlType: SQLTYPE_DELETE}
	op.applyOpts(params)
	sqlData := &SqlData{}
	getTableName(obj, sqlData)
	parseStructSql(obj, sqlData, op)
	return deleteSqlStr(sqlData)
}

func DeleteSql(obj interface{}, params ...OpOption) bool {
	str := DeleteSqlStr(obj, params...)
	return exec(str) == nil
}
