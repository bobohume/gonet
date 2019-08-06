db 基础类库(简单版orm,还原mysql操作,orm和mysql语法一致,没有多余封装,重度orm,程序员对mysql不友善,还在仍受orm无法设置nil,0,fale,"",来吧兄弟)

    db
        将mysql rowresult反馈回来的变更集,按float32,float64,int,int64,string,blob,time转换
       
     db_test
        db 测试代码
        
         列如  
           SimplePlayerData struct{
                AccountId int64 `sql:"primary;name:account_id"`//primary主键,name SQLCOLUMNNAME
                PlayerId int64 `sql:"primary;name:player_id"`
                PlayerName string `sql:"name:player_name"`
                Level int `sql:"name:level"`
                Sex	  int `sql:"name:sex"`
                Gold  int `sql:"name:gold"`
                DrawGold int `sql:"name:draw_gold"`
                Vip int `sql:"name:vip"`
                LastLogoutTime int64 `sql:"datetime;nameg:last_logout_time"`//datetime将int64转化为sql时间类型
                LastLoginTime int64	`sql:"datetime;name:last_login_time"`
                PLayerBlob *message.PlayerData	`sql:"blob;name:plaeyr_blob"`//blob对应mysql的json或者text
                PLayerBlobJson *AA	`sql:"json;name:plaeyr_blob_json"`
            } 
            
       		
    deleteSql
        删除类orm
        t.Log(db.DeleteSql(data, "tbl_player"))//更具主键拼接
        t.Log(db.DeleteSqlEx(data, "tbl_player", "player_id"))//更具主键有选择拼接
        
    insertsql
        插入orm
        	t.Log(db.InsertSql(data, "tbl_player"))
        	t.Log(db.InsertSqlEx(data, "tbl_player", "account_id", "gold", "plaeyr_blob", "plaeyr_blob_json"))
               
    loadsql
        读取orm
        	t.Log(db.LoadSql(data, "tbl_player", "where player_id = 0"))
        	t.Log(db.LoadSqlEx(data, "tbl_player", "where player_id = 0", "account_id", "gold", "plaeyr_blob", "plaeyr_blob_json"))
        	
    updatesql
        更新orm
      	t.Log(db.UpdateSql(data, "tbl_player"))
      	t.Log(db.UpdateSqlEx(data, "tbl_player", "gold", "plaeyr_blob", "plaeyr_blob_json"))
      	
     loadobjsql
        读取orm,读取结果绑定struct
        
        

       