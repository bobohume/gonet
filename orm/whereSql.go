package orm

import (
	"strings"
)

func whereSqlStr(sqlData *SqlData) string {
	key := sqlData.Key
	index := strings.LastIndex(key, ",")
	if index != -1 {
		key = key[:index]
	}
	key = strings.Replace(key, ",", " and ", -1)
	return key
}

func WhereSql(obj interface{}, params ...OpOption) string {
	params = append(params, WithOutWhere())
	op := &Op{sqlType: SQLTYPE_WHERE}
	op.applyOpts(params)
	sqlData := &SqlData{}
	getTableName(obj, sqlData)
	parseStructSql(obj, sqlData, op)
	return whereSqlStr(sqlData)
}
