package db

import (
	"base"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
	"time"
)


type(
	Datetime int64

	Row struct {
		m_Resut map[string] string
	}

	IRow interface {
		init()
		Set(key, val string)
		Get(key string) string
		String(key string) string
		Int(key string) int
		Int64(key string) int64
		Float32(key string) float32
		Float64(key string) float64
		Bool(key string) bool
		Time(key string) int64
		Obj(obj interface{}) bool
	}

	Rows struct {
		m_Rows []*Row
		m_posRow int
	}

	IRows interface {
		Next() bool
		Row() *Row
		Obj(obj interface{}) bool
	}
)

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

func OpenDB(svr string, usr string, pwd string, db string) *sql.DB {
	sqlstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", usr, pwd, svr, db)
	mydb, err := sql.Open("mysql", sqlstr)
	base.ChechErr(err)
	return mydb
}

func getSqlName(sf reflect.StructField) string{
	tagMap := base.ParseTag(sf, "sql")
	if name, exist := tagMap["name"];exist{
		return name
	}

	return strings.ToLower(sf.Name)
}

func isPrimary(sf reflect.StructField) bool{
	tagMap := base.ParseTag(sf, "sql")
	if _, exist := tagMap["primary"];exist{
		return true
	}

	return false
}

func isDatetime(sf reflect.StructField) bool{
	tagMap := base.ParseTag(sf, "sql")
	if _, exist := tagMap["datetime"];exist{
		return true
	}

	return false
}

func isIgnore(sf reflect.StructField) bool{
	tagMap := base.ParseTag(sf, "sql")
	if _, exist := tagMap["-"];exist{
		return true
	}

	return false
}

func (this *Row) init() {
	this.m_Resut = make(map[string] string)
}

func (this *Row) Set(key, val string){
	this.m_Resut[key] = val
}

func (this *Row) Get(key string) string{
	//key = strings.ToLower(key)
	v, exist := this.m_Resut[key]
	if exist{
		return v
	}

	return ""
}

func (this *Row) String(key string) string{
	return this.Get(key)
}

func (this *Row) Int(key string) int{
	n, _ := strconv.Atoi(this.Get(key))
	return n
}

func (this *Row) Int64(key string) int64{
	n, _ := strconv.ParseInt(this.Get(key), 0, 64)
	return n
}

func (this *Row) Float32(key string) float32{
	n, _ := strconv.ParseFloat(this.Get(key), 32)
	return float32(n)
}

func (this *Row) Float64(key string) float64{
	n, _ := strconv.ParseFloat(this.Get(key), 64)
	return n
}

func (this *Row) Bool(key string) bool{
	n, _ := strconv.ParseBool(this.Get(key))
	return n
}

func (this *Row) Time(key string) int64{
	return base.GetDBTime(this.Get(key)).Unix()
}

func (this *Row) Obj(obj interface{}) bool{
	return LoadObjSql(obj, this)
}

func (this *Rows) init(){
	this.m_posRow = 0
}

func (this *Rows) Next() bool{
	if this.m_posRow < len(this.m_Rows){
		this.m_posRow++
		return true
	}
	return false
}

func (this *Rows) Row() *Row{
	nPos := this.m_posRow-1
	if nPos >= 0 && nPos < len(this.m_Rows){
		return this.m_Rows[nPos]
	}

	return NewRow()
}

func (this *Rows) Obj(obj interface{}) bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("rows load obj", err)
		}
	}()

	r := reflect.Indirect(reflect.ValueOf(obj))
	isPtr := false
	if kind := r.Kind(); kind == reflect.Slice {
		rType := r.Type().Elem()
		if rType.Kind() == reflect.Ptr {
			isPtr = true
			rType = rType.Elem()
		}
		for this.Next(){
			elem := reflect.New(rType).Elem()
			LoadObjSql(elem.Addr().Interface(), this.Row())
			if isPtr{
				r.Set(reflect.Append(r, elem.Addr()))
			}else{
				r.Set(reflect.Append(r, elem))
			}
		}
	}
	return true
}

func NewRow() *Row{
	row := &Row{}
	row.init()
	return row
}

func Query(rows *sql.Rows) *Rows{
	rs := &Rows{}
	rs.init()
	if rows != nil{
		cloumns, err := rows.Columns()
		cloumnsLen := len(cloumns)
		if err == nil && cloumnsLen > 0{
			for rows.Next(){
				r := NewRow()
				value := make([]*string, cloumnsLen)
				value1 := make([]interface{}, cloumnsLen)
				for i, _ := range value{
					value[i] = new(string)
					value1[i] = value[i]
				}
				rows.Scan(value1...)
				for i, v := range value{
					r.m_Resut[cloumns[i]] = *v
				}
				rs.m_Rows = append(rs.m_Rows, r)
			}
		}
		rows.Close()
	}
	return rs
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
