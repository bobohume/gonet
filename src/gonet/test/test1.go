package main

//easyjson:json
type Coinmarketcap struct {
	DATA map [string]Coin     `json:"data" bson:"data"`
	MD   [100] Metadata `json:"metadata" bson:"metadata"`
}

//easyjson:json
type Metadata struct {
	Timestamp int   `json:"timestamp" bson:"timestamp"`
	Num       int   `json:"num_cryptocurrencies" bson:"num_cryptocurrencies"`
	Error     error `json:"error" bson:"error"`
}

//easyjson:json
type Coin struct {
	Id          [100]int     `json:"id" bson:"id"`
	Name        [100]string  `json:"name" bson:"name"`
	Symbol      [100]string  `json:"symbol" bson:"symbol"`
	Website     [100]string  `json:"website_slug" bson:"website_slug"`
	Rank        [100]int     `json:"rank" bson:"rank"`
	Circulating [100]float32 `json:"circulating_supply" bson:"circulating_supply"`
	Total       [100]float32 `json:"total_supply" bson:"total_supply"`
	QUotes      [100]Quotes  `json:"quotes" bson:"quotes"`
	Update      [100]int     `json:"last_updated" bson:"last_updated"`
}

//easyjson:json
type Data struct {
	COin Coin `json:"coin"`
}

//easyjson:json
type Quotes struct {
	USD Usd `json:"USD" bson:"USD"`
}

//easyjson:json
type Usd struct {
	Price      float64 `json:"price"	bson:"price"`
	Volume     float64 `json:"volume_24h" bson:"volume_24h"`
	MarketCap  float32 `json:"market_cap" bson:"market_cap"`
	Percent1h  float32 `json:"percent_change_1h" bson:"percent_change_1h"`
	Percent24h float32 `json:"percent_change_24h" bson:"percent_change_24h"`
	Percent7d  float32 `json:"percent_change_7d" bson:"percent_change_7d"`
}

