package model

type(
	TopRank struct {
		table 	 string	`sql:"table;name:tbl_toprank"`
		Id       int64  `sql:"primary;name:id" json:"id"`
		Type     int8   `sql:"primary;name:type" json:"type"`
		Name     string `sql:"name:name" json:"name"`
		Score    int    `sql:"name:score" json:"score"`
		Value    [2]int `sql:"name:value" json:"value"`
		LastTime int64  `sql:"datetime;name:last_time" json:"last_time"`
	}
)
