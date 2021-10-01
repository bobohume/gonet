package model

type(
	MailItem struct{
		table string	`sql:"table;name:tbl_mail"`
		Id int64`sql:"primary;name:id"`
		Sender int64 `sql:"name:sender"`
		SenderName string `sql:"name:sender_name"`
		Recver int64 `sql:"name:recver"`
		RecverName string `sql:"name:recver_name"`
		Money int `sql:"name:money"`
		ItemId int `sql:"name:item_id"`
		ItemCount int `sql:"name:item_count"`
		IsRead int8 `sql:"name:is_read"`
		IsSystem int8 `sql:"name:is_system"`
		RecvFlag int8 `sql:"name:recv_flag"`
		Title string `sql:"name:title"`
		Content string `sql:"name:content"`
	}
)
