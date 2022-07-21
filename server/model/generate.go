package model

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type (
	moduleInfo struct {
		className  string
		fileName   string
		moduleName string
	}

	fileInfo struct {
		obj        interface{}
		moduleList []moduleInfo
		dbKey      string
		dbKeyType  string
	}
)

var (
	FILE_MAP = map[string]fileInfo{
		"Player": fileInfo{
			obj: PlayerData{},
			moduleList: []moduleInfo{
				moduleInfo{"Player", "./../../server/game/player/PlayerData.go", "game"},
				moduleInfo{"Player", "./../../server/db/PlayerData.go", "db"},
			},
			dbKey:     `PlayerId`,
			dbKeyType: `int64`,
		},
	}

	//包头
	FILE_GENERATE_HEAD_MAP = map[string]string{
		"game": `package player

import(
	"gonet/base"
	"gonet/common/cluster"
	"gonet/rpc"
)

// 自动生成代码
`,

		"db": `package db

import(
	"context"
	"gonet/base"
    "gonet/orm"
    "gonet/server/model"
)

// 自动生成代码
`,
	}
	//包体
	FILE_GENERATE_BODY_MAP = map[string][]string{
		"game": []string{
			`
func (this *{ClassName}) Save{MemberType}(){
	this.{MemberName}.Dirty = true
}

func (this *{ClassName}) __Save{MemberType}DB(){
	if this.{MemberName}.Dirty{
    	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, Id:this.MailBox.Id}, "PlayerMgr.Save{MemberType}", this.{MemberType}.{DbKeyName}, this.{MemberName})
		this.{MemberName}.Dirty = false
    	base.LOG.Printf("玩家[%d] Save{MemberType}", this.MailBox.Id)
	}
}
`,
		},

		"db": []string{
			`
func (this *{ClassName}) __Load{MemberType}DB({DbKeyName} {DbKeyTypeName}) error{
    data := &model.{MemberType}{{DbKeyName}:{DbKeyName}}
    rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.{MemberName}, rs.Row())
    }
	return err
}
`,
			`
/*
func (this *{ClassName}Mgr) Save{MemberType}(ctx context.Context, playerId int64, data model.{MemberType}){
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] Save{MemberType}", playerId)
}
*/
func (this *{ClassName}) __Save{MemberType}DB(){
	if this.{MemberName}.Dirty{
    	orm.DB.Exec(orm.SaveSql(this.{MemberName}))
		this.{MemberName}.Dirty = false
	}
}

func (this *{ClassName}) __Save{MemberType}(data model.{MemberType}){
    this.{MemberName} = data
	this.{MemberName}.Dirty = true
    base.LOG.Printf("玩家[%d] Save{MemberType}", this.{MemberType}.{DbKeyName})
}

func (this *{ClassName}Mgr) Save{MemberType}(ctx context.Context, playerId int64, data model.{MemberType}){
    player := this.GetPlayer(playerId)
	if player != nil{
        player.__Save{MemberType}(data)
	}
}
`,
		},
	}
)

func Generate(name string) {
	fileInfo := FILE_MAP[name]
	classType := reflect.TypeOf(fileInfo.obj)
	for _, moduleInfo := range fileInfo.moduleList {
		className := moduleInfo.className
		stream := bytes.NewBuffer([]byte{})
		stream.WriteString(FILE_GENERATE_HEAD_MAP[moduleInfo.moduleName])
		memberNameList := []string{}
		for i := 0; i < classType.NumField(); i++ {
			sf := classType.Field(i)
			memberType := sf.Type.Name()
			memberName := sf.Name
			memberNameList = append(memberNameList, memberName)
			for _, funcStr := range FILE_GENERATE_BODY_MAP[moduleInfo.moduleName] {
				str := funcStr
				str = strings.Replace(str, "{ClassName}", className, -1)
				str = strings.Replace(str, "{MemberType}", memberType, -1)
				str = strings.Replace(str, "{MemberName}", memberName, -1)
				str = strings.Replace(str, "{DbKeyName}", fileInfo.dbKey, -1)
				str = strings.Replace(str, "{DbKeyTypeName}", fileInfo.dbKeyType, -1)
				stream.WriteString(str)
			}
		}

		if moduleInfo.moduleName == "db" {
			// func LoadDB
			stream.WriteString(fmt.Sprintf("\nfunc (this *%s) Load%sDB(%s %s) error{\n", className, name, fileInfo.dbKey, fileInfo.dbKeyType))
			stream.WriteString(fmt.Sprintf("    this.Init(%s)\n", fileInfo.dbKey))
			for _, v := range memberNameList {
				stream.WriteString(fmt.Sprintf("    if err := this.__Load%sDB(%s); err != nil{\n", v, fileInfo.dbKey))
				stream.WriteString(fmt.Sprintf(`        base.LOG.Printf("__Load%sDB() error")`+"\n", v))
				stream.WriteString(fmt.Sprintf("        return err \n"))
				stream.WriteString(fmt.Sprintf("    }\n"))
			}
			stream.WriteString("    return nil\n")
			stream.WriteString("}\n\n")

			// func SaveDB
			stream.WriteString(fmt.Sprintf("\nfunc (this *%s) Save%sDB(){\n", className, name))
			for _, v := range memberNameList {
				stream.WriteString(fmt.Sprintf("    this.__Save%sDB()\n", v))
			}
			stream.WriteString("}\n\n")
		} else if moduleInfo.moduleName == "game" {
			// func SaveDB
			stream.WriteString(fmt.Sprintf("\nfunc (this *%s) Save%sDB(){\n", className, name))
			for _, v := range memberNameList {
				stream.WriteString(fmt.Sprintf("    this.__Save%sDB()\n", v))
			}
			stream.WriteString("}\n\n")
		}

		file, err := os.Create(moduleInfo.fileName)
		if err == nil {
			file.Write(stream.Bytes())
			file.Close()
		}
		fmt.Println(err)
	}
}
