package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"time"
	"reflect"
	"strings"
)


type(
	Datetime int64
)

func chechErr(err error) {
	if err == nil {
		return
	}
	log.Fatalf("错误：%s\n", err.Error())
}

// BinaryData的类型是blob/longblob
func addData(db *sql.DB, data []byte) error {
	_, err := db.Exec("INSERT INTO MyTable(BinaryData) VALUES(?)", data)
	return err
}

//---获取datetime时间
func  GetDBTimeString(t int64)string{
	tm := time.Unix(t, 0)
	return  tm.Format("2006-01-02 15:04:05")
}

func GetDBTime(strTime string) *time.Time{
	DefaultTimeLoc := time.Local
	loginTime, err := time.ParseInLocation("2006-01-02 15:04:05", strTime, DefaultTimeLoc)
	chechErr(err)
	return &loginTime
}

func OpenDB(svr string, usr string, pwd string, db string) *sql.DB {
	sqlstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", usr, pwd, svr, db)
	mydb, err := sql.Open("mysql", sqlstr)
	chechErr(err)
	return mydb
}

func isPrimary(classField reflect.StructField) bool{
	tempstr := fmt.Sprintf("%v", classField)
	if strings.Index(tempstr, "primary") != -1{
		return true
	}
	return false
}

func isDatetime(classField reflect.StructField) bool{
	tempstr := fmt.Sprintf("%v", classField)
	if strings.Index(tempstr, "datetime") != -1{
		return true
	}
	return false
}
//example
/*db := db.OpenDB("localhost:3306", "root", "123456", "test")
nlen := 0
for nlen < 10 {
	go func() {
		for {
			db.Exec("INSERT INTO tbl_car(carnum) values(?)", 2)
		}
	}()
	nlen++
}*/
