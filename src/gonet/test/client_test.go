package main_test

import (
	"encoding/json"
	"github.com/json-iterator/go"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type Coinmarketcap struct {
	DATA map [string]Coin     `json:"data" bson:"data"`
	MD   Metadata `json:"metadata" bson:"metadata"`
}

type Metadata struct {
	Timestamp int   `json:"timestamp" bson:"timestamp"`
	Num       int   `json:"num_cryptocurrencies" bson:"num_cryptocurrencies"`
	Error     error `json:"error" bson:"error"`
}

type Coin struct {
	Id          int     `json:"id" bson:"id"`
	Name        string  `json:"name" bson:"name"`
	Symbol      string  `json:"symbol" bson:"symbol"`
	Website     string  `json:"website_slug" bson:"website_slug"`
	Rank        int     `json:"rank" bson:"rank"`
	Circulating float32 `json:"circulating_supply" bson:"circulating_supply"`
	Total       float32 `json:"total_supply" bson:"total_supply"`
	QUotes      Quotes  `json:"quotes" bson:"quotes"`
	Update      int     `json:"last_updated" bson:"last_updated"`
}

type Data struct {
	COin Coin `json:"coin"`
}

type Quotes struct {
	USD Usd `json:"USD" bson:"USD"`
}

type Usd struct {
	Price      float64 `json:"price"	bson:"price"`
	Volume     float64 `json:"volume_24h" bson:"volume_24h"`
	MarketCap  float32 `json:"market_cap" bson:"market_cap"`
	Percent1h  float32 `json:"percent_change_1h" bson:"percent_change_1h"`
	Percent24h float32 `json:"percent_change_24h" bson:"percent_change_24h"`
	Percent7d  float32 `json:"percent_change_7d" bson:"percent_change_7d"`
}

var(
	data1 = []byte{}
	data = []byte(`{
    "attention": "WARNING: This API is now deprecated and will be taken offline soon.  Please switch to the new CoinMarketCap API to avoid interruptions in service. (https://pro.coinmarketcap.com/migrate/)", 
    "data": {
        "1": {
            "id": 1, 
            "name": "Bitcoin", 
            "symbol": "BTC", 
            "website_slug": "bitcoin", 
            "rank": 1, 
            "circulating_supply": 17605287.0, 
            "total_supply": 17605287.0, 
            "max_supply": 21000000.0, 
            "quotes": {
                "USD": {
                    "price": 4037.53370383, 
                    "volume_24h": 9474994090.1506, 
                    "market_cap": 71081939628.0, 
                    "percent_change_1h": -0.21, 
                    "percent_change_24h": 0.08, 
                    "percent_change_7d": -0.37
                }
            }, 
            "last_updated": 1553335765
        }, 
        "1027": {
            "id": 1027, 
            "name": "Ethereum", 
            "symbol": "ETH", 
            "website_slug": "ethereum", 
            "rank": 2, 
            "circulating_supply": 105353289.0, 
            "total_supply": 105353289.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 139.09544543, 
                    "volume_24h": 4495596246.78627, 
                    "market_cap": 14654162661.0, 
                    "percent_change_1h": -0.09, 
                    "percent_change_24h": 1.12, 
                    "percent_change_7d": -2.7
                }
            }, 
            "last_updated": 1553335760
        }, 
        "52": {
            "id": 52, 
            "name": "XRP", 
            "symbol": "XRP", 
            "website_slug": "ripple", 
            "rank": 3, 
            "circulating_supply": 41666017553.0, 
            "total_supply": 99991672219.0, 
            "max_supply": 100000000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.3128065494, 
                    "volume_24h": 673574840.84748, 
                    "market_cap": 13033403180.0, 
                    "percent_change_1h": -0.07, 
                    "percent_change_24h": 0.12, 
                    "percent_change_7d": -2.38
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2": {
            "id": 2, 
            "name": "Litecoin", 
            "symbol": "LTC", 
            "website_slug": "litecoin", 
            "rank": 4, 
            "circulating_supply": 61007136.0, 
            "total_supply": 61007136.0, 
            "max_supply": 84000000.0, 
            "quotes": {
                "USD": {
                    "price": 61.5552375765, 
                    "volume_24h": 1874989766.86926, 
                    "market_cap": 3755308780.0, 
                    "percent_change_1h": -0.08, 
                    "percent_change_24h": 3.73, 
                    "percent_change_7d": -0.04
                }
            }, 
            "last_updated": 1553335746
        }, 
        "1765": {
            "id": 1765, 
            "name": "EOS", 
            "symbol": "EOS", 
            "website_slug": "eos", 
            "rank": 5, 
            "circulating_supply": 906245118.0, 
            "total_supply": 1006245120.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 3.7039303795, 
                    "volume_24h": 1435319703.15984, 
                    "market_cap": 3356668822.0, 
                    "percent_change_1h": -0.25, 
                    "percent_change_24h": 1.41, 
                    "percent_change_7d": -3.22
                }
            }, 
            "last_updated": 1553335744
        }, 
        "1831": {
            "id": 1831, 
            "name": "Bitcoin Cash", 
            "symbol": "BCH", 
            "website_slug": "bitcoin-cash", 
            "rank": 6, 
            "circulating_supply": 17687863.0, 
            "total_supply": 17687863.0, 
            "max_supply": 21000000.0, 
            "quotes": {
                "USD": {
                    "price": 168.5644602, 
                    "volume_24h": 506174965.194053, 
                    "market_cap": 2981544994.0, 
                    "percent_change_1h": -0.56, 
                    "percent_change_24h": 8.55, 
                    "percent_change_7d": 7.52
                }
            }, 
            "last_updated": 1553335745
        }, 
        "1839": {
            "id": 1839, 
            "name": "Binance Coin", 
            "symbol": "BNB", 
            "website_slug": "binance-coin", 
            "rank": 7, 
            "circulating_supply": 141175490.0, 
            "total_supply": 189175490.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 15.383564469, 
                    "volume_24h": 153864421.77578, 
                    "market_cap": 2171782256.0, 
                    "percent_change_1h": -0.24, 
                    "percent_change_24h": 3.88, 
                    "percent_change_7d": -1.35
                }
            }, 
            "last_updated": 1553335744
        }, 
        "512": {
            "id": 512, 
            "name": "Stellar", 
            "symbol": "XLM", 
            "website_slug": "stellar", 
            "rank": 8, 
            "circulating_supply": 19225306519.0, 
            "total_supply": 104822379135.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.1077841468, 
                    "volume_24h": 187178349.511135, 
                    "market_cap": 2072183259.0, 
                    "percent_change_1h": -0.72, 
                    "percent_change_24h": 0.12, 
                    "percent_change_7d": -1.81
                }
            }, 
            "last_updated": 1553335743
        }, 
        "825": {
            "id": 825, 
            "name": "Tether", 
            "symbol": "USDT", 
            "website_slug": "tether", 
            "rank": 9, 
            "circulating_supply": 2020708392.0, 
            "total_supply": 2580057493.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.0085441309, 
                    "volume_24h": 8341592198.5792, 
                    "market_cap": 2037973589.0, 
                    "percent_change_1h": 0.01, 
                    "percent_change_24h": -0.34, 
                    "percent_change_7d": -0.31
                }
            }, 
            "last_updated": 1553335752
        }, 
        "2010": {
            "id": 2010, 
            "name": "Cardano", 
            "symbol": "ADA", 
            "website_slug": "cardano", 
            "rank": 10, 
            "circulating_supply": 25927070538.0, 
            "total_supply": 31112483745.0, 
            "max_supply": 45000000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.059891365, 
                    "volume_24h": 126265663.755353, 
                    "market_cap": 1552807644.0, 
                    "percent_change_1h": -0.42, 
                    "percent_change_24h": 5.25, 
                    "percent_change_7d": 17.64
                }
            }, 
            "last_updated": 1553335744
        }, 
        "1958": {
            "id": 1958, 
            "name": "TRON", 
            "symbol": "TRX", 
            "website_slug": "tron", 
            "rank": 11, 
            "circulating_supply": 66682072191.0, 
            "total_supply": 99281283754.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0230691962, 
                    "volume_24h": 186970746.733909, 
                    "market_cap": 1538301803.0, 
                    "percent_change_1h": 0.44, 
                    "percent_change_24h": 1.75, 
                    "percent_change_7d": -1.73
                }
            }, 
            "last_updated": 1553335745
        }, 
        "3602": {
            "id": 3602, 
            "name": "Bitcoin SV", 
            "symbol": "BSV", 
            "website_slug": "bitcoin-sv", 
            "rank": 12, 
            "circulating_supply": 17670348.0, 
            "total_supply": 17670348.0, 
            "max_supply": 21000000.0, 
            "quotes": {
                "USD": {
                    "price": 67.8835937183, 
                    "volume_24h": 110069553.093655, 
                    "market_cap": 1199526724.0, 
                    "percent_change_1h": 0.1, 
                    "percent_change_24h": 2.07, 
                    "percent_change_7d": -2.81
                }
            }, 
            "last_updated": 1553335749
        }, 
        "328": {
            "id": 328, 
            "name": "Monero", 
            "symbol": "XMR", 
            "website_slug": "monero", 
            "rank": 13, 
            "circulating_supply": 16871262.0, 
            "total_supply": 16871262.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 54.1709162834, 
                    "volume_24h": 85422504.2165936, 
                    "market_cap": 913931696.0, 
                    "percent_change_1h": 0.1, 
                    "percent_change_24h": 1.06, 
                    "percent_change_7d": -1.13
                }
            }, 
            "last_updated": 1553335742
        }, 
        "1720": {
            "id": 1720, 
            "name": "IOTA", 
            "symbol": "MIOTA", 
            "website_slug": "iota", 
            "rank": 14, 
            "circulating_supply": 2779530283.0, 
            "total_supply": 2779530283.0, 
            "max_supply": 2779530283.0, 
            "quotes": {
                "USD": {
                    "price": 0.3171575084, 
                    "volume_24h": 17970797.2727939, 
                    "market_cap": 881548899.0, 
                    "percent_change_1h": -0.26, 
                    "percent_change_24h": 1.66, 
                    "percent_change_7d": 3.34
                }
            }, 
            "last_updated": 1553335743
        }, 
        "131": {
            "id": 131, 
            "name": "Dash", 
            "symbol": "DASH", 
            "website_slug": "dash", 
            "rank": 15, 
            "circulating_supply": 8704127.0, 
            "total_supply": 8704127.0, 
            "max_supply": 18900000.0, 
            "quotes": {
                "USD": {
                    "price": 92.5974925131, 
                    "volume_24h": 282997077.98902, 
                    "market_cap": 805980344.0, 
                    "percent_change_1h": 0.01, 
                    "percent_change_24h": 0.97, 
                    "percent_change_7d": -0.15
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1518": {
            "id": 1518, 
            "name": "Maker", 
            "symbol": "MKR", 
            "website_slug": "maker", 
            "rank": 16, 
            "circulating_supply": 1000000.0, 
            "total_supply": 1000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 733.623857902, 
                    "volume_24h": 2061271.10662478, 
                    "market_cap": 733623858.0, 
                    "percent_change_1h": -0.4, 
                    "percent_change_24h": 2.21, 
                    "percent_change_7d": 4.97
                }
            }, 
            "last_updated": 1553335743
        }, 
        "2566": {
            "id": 2566, 
            "name": "Ontology", 
            "symbol": "ONT", 
            "website_slug": "ontology", 
            "rank": 17, 
            "circulating_supply": 494823234.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.31284924, 
                    "volume_24h": 107977655.019428, 
                    "market_cap": 649628307.0, 
                    "percent_change_1h": 0.09, 
                    "percent_change_24h": 1.16, 
                    "percent_change_7d": 23.21
                }
            }, 
            "last_updated": 1553335746
        }, 
        "1376": {
            "id": 1376, 
            "name": "NEO", 
            "symbol": "NEO", 
            "website_slug": "neo", 
            "rank": 18, 
            "circulating_supply": 65000000.0, 
            "total_supply": 100000000.0, 
            "max_supply": 100000000.0, 
            "quotes": {
                "USD": {
                    "price": 9.3694766414, 
                    "volume_24h": 266877801.614936, 
                    "market_cap": 609015982.0, 
                    "percent_change_1h": -0.1, 
                    "percent_change_24h": 1.08, 
                    "percent_change_7d": -2.74
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1321": {
            "id": 1321, 
            "name": "Ethereum Classic", 
            "symbol": "ETC", 
            "website_slug": "ethereum-classic", 
            "rank": 19, 
            "circulating_supply": 109172414.0, 
            "total_supply": 109172414.0, 
            "max_supply": 210000000.0, 
            "quotes": {
                "USD": {
                    "price": 4.9112832392, 
                    "volume_24h": 229610995.433755, 
                    "market_cap": 536176647.0, 
                    "percent_change_1h": 0.07, 
                    "percent_change_24h": -0.39, 
                    "percent_change_7d": 7.9
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2011": {
            "id": 2011, 
            "name": "Tezos", 
            "symbol": "XTZ", 
            "website_slug": "tezos", 
            "rank": 20, 
            "circulating_supply": 665508583.0, 
            "total_supply": 788216816.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.7536561596, 
                    "volume_24h": 4136599.38060336, 
                    "market_cap": 501564643.0, 
                    "percent_change_1h": -0.93, 
                    "percent_change_24h": -0.08, 
                    "percent_change_7d": 63.75
                }
            }, 
            "last_updated": 1553335744
        }, 
        "873": {
            "id": 873, 
            "name": "NEM", 
            "symbol": "XEM", 
            "website_slug": "nem", 
            "rank": 21, 
            "circulating_supply": 8999999999.0, 
            "total_supply": 8999999999.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.050712143, 
                    "volume_24h": 12689091.9451282, 
                    "market_cap": 456409287.0, 
                    "percent_change_1h": -0.03, 
                    "percent_change_24h": 2.7, 
                    "percent_change_7d": 0.72
                }
            }, 
            "last_updated": 1553335742
        }, 
        "1437": {
            "id": 1437, 
            "name": "Zcash", 
            "symbol": "ZEC", 
            "website_slug": "zcash", 
            "rank": 22, 
            "circulating_supply": 6156531.0, 
            "total_supply": 6156531.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 56.4950739664, 
                    "volume_24h": 171307548.906149, 
                    "market_cap": 347813688.0, 
                    "percent_change_1h": -0.56, 
                    "percent_change_24h": 0.8, 
                    "percent_change_7d": 2.68
                }
            }, 
            "last_updated": 1553335744
        }, 
        "3077": {
            "id": 3077, 
            "name": "VeChain", 
            "symbol": "VET", 
            "website_slug": "vechain", 
            "rank": 23, 
            "circulating_supply": 55454734800.0, 
            "total_supply": 86712634466.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0060193822, 
                    "volume_24h": 19623211.2163555, 
                    "market_cap": 333803242.0, 
                    "percent_change_1h": -0.73, 
                    "percent_change_24h": 3.26, 
                    "percent_change_7d": 11.81
                }
            }, 
            "last_updated": 1553335747
        }, 
        "3635": {
            "id": 3635, 
            "name": "Crypto.com Chain", 
            "symbol": "CRO", 
            "website_slug": "crypto-com-chain", 
            "rank": 24, 
            "circulating_supply": 4200913242.0, 
            "total_supply": 100000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0673041719, 
                    "volume_24h": 742979.020422371, 
                    "market_cap": 282738987.0, 
                    "percent_change_1h": 0.45, 
                    "percent_change_24h": 2.49, 
                    "percent_change_7d": -27.22
                }
            }, 
            "last_updated": 1553335749
        }, 
        "1274": {
            "id": 1274, 
            "name": "Waves", 
            "symbol": "WAVES", 
            "website_slug": "waves", 
            "rank": 25, 
            "circulating_supply": 100000000.0, 
            "total_supply": 100000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 2.7944801605, 
                    "volume_24h": 7344068.20084145, 
                    "market_cap": 279448016.0, 
                    "percent_change_1h": 0.07, 
                    "percent_change_24h": 1.13, 
                    "percent_change_7d": -0.61
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1697": {
            "id": 1697, 
            "name": "Basic Attention Token", 
            "symbol": "BAT", 
            "website_slug": "basic-attention-token", 
            "rank": 26, 
            "circulating_supply": 1244066783.0, 
            "total_supply": 1500000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.201510088, 
                    "volume_24h": 13049756.846002, 
                    "market_cap": 250692007.0, 
                    "percent_change_1h": 0.51, 
                    "percent_change_24h": 1.82, 
                    "percent_change_7d": 1.03
                }
            }, 
            "last_updated": 1553335743
        }, 
        "3408": {
            "id": 3408, 
            "name": "USD Coin", 
            "symbol": "USDC", 
            "website_slug": "usd-coin", 
            "rank": 27, 
            "circulating_supply": 244698122.0, 
            "total_supply": 244963802.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.0110731915, 
                    "volume_24h": 21207372.796982, 
                    "market_cap": 247407712.0, 
                    "percent_change_1h": -0.1, 
                    "percent_change_24h": -0.4, 
                    "percent_change_7d": -0.4
                }
            }, 
            "last_updated": 1553335749
        }, 
        "74": {
            "id": 74, 
            "name": "Dogecoin", 
            "symbol": "DOGE", 
            "website_slug": "dogecoin", 
            "rank": 28, 
            "circulating_supply": 118765925285.0, 
            "total_supply": 118765925285.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0020605903, 
                    "volume_24h": 17092954.3507101, 
                    "market_cap": 244727915.0, 
                    "percent_change_1h": 0.02, 
                    "percent_change_24h": 0.09, 
                    "percent_change_7d": -0.84
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1808": {
            "id": 1808, 
            "name": "OmiseGO", 
            "symbol": "OMG", 
            "website_slug": "omisego", 
            "rank": 29, 
            "circulating_supply": 140245398.0, 
            "total_supply": 140245398.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.7306515694, 
                    "volume_24h": 81249431.6677763, 
                    "market_cap": 242715919.0, 
                    "percent_change_1h": 0.76, 
                    "percent_change_24h": 9.57, 
                    "percent_change_7d": 12.55
                }
            }, 
            "last_updated": 1553335744
        }, 
        "1684": {
            "id": 1684, 
            "name": "Qtum", 
            "symbol": "QTUM", 
            "website_slug": "qtum", 
            "rank": 30, 
            "circulating_supply": 89343636.0, 
            "total_supply": 101343636.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 2.6679857068, 
                    "volume_24h": 159825831.066343, 
                    "market_cap": 238367544.0, 
                    "percent_change_1h": -0.41, 
                    "percent_change_24h": 5.4, 
                    "percent_change_7d": 5.18
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2083": {
            "id": 2083, 
            "name": "Bitcoin Gold", 
            "symbol": "BTG", 
            "website_slug": "bitcoin-gold", 
            "rank": 31, 
            "circulating_supply": 17413924.0, 
            "total_supply": 17513924.0, 
            "max_supply": 21000000.0, 
            "quotes": {
                "USD": {
                    "price": 13.2674731325, 
                    "volume_24h": 10313924.4791074, 
                    "market_cap": 231038763.0, 
                    "percent_change_1h": -0.05, 
                    "percent_change_24h": 1.55, 
                    "percent_change_7d": -4.51
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2563": {
            "id": 2563, 
            "name": "TrueUSD", 
            "symbol": "TUSD", 
            "website_slug": "trueusd", 
            "rank": 32, 
            "circulating_supply": 203095205.0, 
            "total_supply": 203095205.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.0159854777, 
                    "volume_24h": 39007521.9259149, 
                    "market_cap": 206341779.0, 
                    "percent_change_1h": -0.08, 
                    "percent_change_24h": 0.01, 
                    "percent_change_7d": -0.5
                }
            }, 
            "last_updated": 1553335746
        }, 
        "1168": {
            "id": 1168, 
            "name": "Decred", 
            "symbol": "DCR", 
            "website_slug": "decred", 
            "rank": 33, 
            "circulating_supply": 9509512.0, 
            "total_supply": 9509512.0, 
            "max_supply": 21000000.0, 
            "quotes": {
                "USD": {
                    "price": 19.0450081632, 
                    "volume_24h": 1757260.93625635, 
                    "market_cap": 181108738.0, 
                    "percent_change_1h": 0.05, 
                    "percent_change_24h": 3.22, 
                    "percent_change_7d": -2.95
                }
            }, 
            "last_updated": 1553335742
        }, 
        "1214": {
            "id": 1214, 
            "name": "Lisk", 
            "symbol": "LSK", 
            "website_slug": "lisk", 
            "rank": 34, 
            "circulating_supply": 115565555.0, 
            "total_supply": 130680685.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.5232404163, 
                    "volume_24h": 4265040.06336163, 
                    "market_cap": 176034123.0, 
                    "percent_change_1h": 0.1, 
                    "percent_change_24h": 1.86, 
                    "percent_change_7d": 0.82
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2469": {
            "id": 2469, 
            "name": "Zilliqa", 
            "symbol": "ZIL", 
            "website_slug": "zilliqa", 
            "rank": 35, 
            "circulating_supply": 8656969775.0, 
            "total_supply": 12599999804.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0200644793, 
                    "volume_24h": 20990774.2797482, 
                    "market_cap": 173697591.0, 
                    "percent_change_1h": -0.36, 
                    "percent_change_24h": 5.21, 
                    "percent_change_7d": 8.42
                }
            }, 
            "last_updated": 1553335746
        }, 
        "1104": {
            "id": 1104, 
            "name": "Augur", 
            "symbol": "REP", 
            "website_slug": "augur", 
            "rank": 36, 
            "circulating_supply": 11000000.0, 
            "total_supply": 11000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 15.5532558208, 
                    "volume_24h": 7744222.68259148, 
                    "market_cap": 171085814.0, 
                    "percent_change_1h": -0.71, 
                    "percent_change_24h": 5.0, 
                    "percent_change_7d": 6.37
                }
            }, 
            "last_updated": 1553335742
        }, 
        "109": {
            "id": 109, 
            "name": "DigiByte", 
            "symbol": "DGB", 
            "website_slug": "digibyte", 
            "rank": 37, 
            "circulating_supply": 11583299045.0, 
            "total_supply": 11583299045.0, 
            "max_supply": 21000000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.0146229349, 
                    "volume_24h": 2766676.6723752, 
                    "market_cap": 169381827.0, 
                    "percent_change_1h": 0.03, 
                    "percent_change_24h": 0.3, 
                    "percent_change_7d": 0.54
                }
            }, 
            "last_updated": 1553335741
        }, 
        "1975": {
            "id": 1975, 
            "name": "Chainlink", 
            "symbol": "LINK", 
            "website_slug": "chainlink", 
            "rank": 38, 
            "circulating_supply": 350000000.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.4651082112, 
                    "volume_24h": 2904655.31885976, 
                    "market_cap": 162787874.0, 
                    "percent_change_1h": 0.02, 
                    "percent_change_24h": 1.62, 
                    "percent_change_7d": -4.1
                }
            }, 
            "last_updated": 1553335744
        }, 
        "1896": {
            "id": 1896, 
            "name": "0x", 
            "symbol": "ZRX", 
            "website_slug": "0x", 
            "rank": 39, 
            "circulating_supply": 586141504.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.2743035322, 
                    "volume_24h": 17615895.7876782, 
                    "market_cap": 160780685.0, 
                    "percent_change_1h": 0.25, 
                    "percent_change_24h": 4.41, 
                    "percent_change_7d": 0.34
                }
            }, 
            "last_updated": 1553335745
        }, 
        "2682": {
            "id": 2682, 
            "name": "Holo", 
            "symbol": "HOT", 
            "website_slug": "holo", 
            "rank": 40, 
            "circulating_supply": 133214575156.0, 
            "total_supply": 177619433541.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0011845153, 
                    "volume_24h": 6544386.02059278, 
                    "market_cap": 157794704.0, 
                    "percent_change_1h": -0.03, 
                    "percent_change_24h": 0.53, 
                    "percent_change_7d": -5.64
                }
            }, 
            "last_updated": 1553335746
        }, 
        "2099": {
            "id": 2099, 
            "name": "ICON", 
            "symbol": "ICX", 
            "website_slug": "icon", 
            "rank": 41, 
            "circulating_supply": 473406688.0, 
            "total_supply": 800460000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.3250990505, 
                    "volume_24h": 11508667.6423704, 
                    "market_cap": 153904065.0, 
                    "percent_change_1h": -0.56, 
                    "percent_change_24h": 1.55, 
                    "percent_change_7d": -4.88
                }
            }, 
            "last_updated": 1553335743
        }, 
        "463": {
            "id": 463, 
            "name": "BitShares", 
            "symbol": "BTS", 
            "website_slug": "bitshares", 
            "rank": 42, 
            "circulating_supply": 2701150000.0, 
            "total_supply": 2701150000.0, 
            "max_supply": 3600570502.0, 
            "quotes": {
                "USD": {
                    "price": 0.0544015989, 
                    "volume_24h": 13451225.9525613, 
                    "market_cap": 146946879.0, 
                    "percent_change_1h": 0.75, 
                    "percent_change_24h": 8.44, 
                    "percent_change_7d": 5.35
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2130": {
            "id": 2130, 
            "name": "Enjin Coin", 
            "symbol": "ENJ", 
            "website_slug": "enjin-coin", 
            "rank": 43, 
            "circulating_supply": 767007985.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.1881209162, 
                    "volume_24h": 14807471.9072077, 
                    "market_cap": 144290245.0, 
                    "percent_change_1h": 0.42, 
                    "percent_change_24h": -0.49, 
                    "percent_change_7d": 12.58
                }
            }, 
            "last_updated": 1553335745
        }, 
        "1230": {
            "id": 1230, 
            "name": "Steem", 
            "symbol": "STEEM", 
            "website_slug": "steem", 
            "rank": 44, 
            "circulating_supply": 307779505.0, 
            "total_supply": 324753599.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.4684412248, 
                    "volume_24h": 1419392.12519906, 
                    "market_cap": 144176608.0, 
                    "percent_change_1h": -0.07, 
                    "percent_change_24h": 0.72, 
                    "percent_change_7d": -0.95
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2577": {
            "id": 2577, 
            "name": "Ravencoin", 
            "symbol": "RVN", 
            "website_slug": "ravencoin", 
            "rank": 45, 
            "circulating_supply": 3163180000.0, 
            "total_supply": 3163180000.0, 
            "max_supply": 21000000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.0450939615, 
                    "volume_24h": 30983828.0963863, 
                    "market_cap": 142640317.0, 
                    "percent_change_1h": 0.29, 
                    "percent_change_24h": 1.44, 
                    "percent_change_7d": 64.14
                }
            }, 
            "last_updated": 1553335745
        }, 
        "372": {
            "id": 372, 
            "name": "Bytecoin", 
            "symbol": "BCN", 
            "website_slug": "bytecoin-bcn", 
            "rank": 46, 
            "circulating_supply": 184066828814.0, 
            "total_supply": 184066828814.0, 
            "max_supply": 184470000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.0007632227, 
                    "volume_24h": 140696.116704439, 
                    "market_cap": 140483981.0, 
                    "percent_change_1h": 1.86, 
                    "percent_change_24h": 1.96, 
                    "percent_change_7d": -5.15
                }
            }, 
            "last_updated": 1553335741
        }, 
        "3718": {
            "id": 3718, 
            "name": "BitTorrent", 
            "symbol": "BTT", 
            "website_slug": "bittorrent", 
            "rank": 47, 
            "circulating_supply": 170421000000.0, 
            "total_supply": 990000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0007883428, 
                    "volume_24h": 16226666.8584553, 
                    "market_cap": 134350165.0, 
                    "percent_change_1h": -0.63, 
                    "percent_change_24h": 0.63, 
                    "percent_change_7d": 3.26
                }
            }, 
            "last_updated": 1553335749
        }, 
        "1567": {
            "id": 1567, 
            "name": "Nano", 
            "symbol": "NANO", 
            "website_slug": "nano", 
            "rank": 48, 
            "circulating_supply": 133248289.0, 
            "total_supply": 133248289.0, 
            "max_supply": 133248290.0, 
            "quotes": {
                "USD": {
                    "price": 0.997453032, 
                    "volume_24h": 1808791.81462783, 
                    "market_cap": 132908910.0, 
                    "percent_change_1h": -0.06, 
                    "percent_change_24h": 1.24, 
                    "percent_change_7d": -2.32
                }
            }, 
            "last_updated": 1553335743
        }, 
        "2222": {
            "id": 2222, 
            "name": "Bitcoin Diamond", 
            "symbol": "BCD", 
            "website_slug": "bitcoin-diamond", 
            "rank": 49, 
            "circulating_supply": 153756875.0, 
            "total_supply": 156756875.0, 
            "max_supply": 210000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.8605547259, 
                    "volume_24h": 1780772.73984948, 
                    "market_cap": 132316205.0, 
                    "percent_change_1h": -0.74, 
                    "percent_change_24h": 1.35, 
                    "percent_change_7d": -5.28
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2502": {
            "id": 2502, 
            "name": "Huobi Token", 
            "symbol": "HT", 
            "website_slug": "huobi-token", 
            "rank": 50, 
            "circulating_supply": 50000200.0, 
            "total_supply": 500000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 2.5127217685, 
                    "volume_24h": 67501569.8592528, 
                    "market_cap": 125636591.0, 
                    "percent_change_1h": 0.51, 
                    "percent_change_24h": -1.55, 
                    "percent_change_7d": 22.93
                }
            }, 
            "last_updated": 1553335745
        }, 
        "1700": {
            "id": 1700, 
            "name": "Aeternity", 
            "symbol": "AE", 
            "website_slug": "aeternity", 
            "rank": 51, 
            "circulating_supply": 254827375.0, 
            "total_supply": 300648318.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.4894113567, 
                    "volume_24h": 45637117.6723041, 
                    "market_cap": 124715411.0, 
                    "percent_change_1h": -0.45, 
                    "percent_change_24h": 4.4, 
                    "percent_change_7d": 3.9
                }
            }, 
            "last_updated": 1553335743
        }, 
        "3330": {
            "id": 3330, 
            "name": "Paxos Standard Token", 
            "symbol": "PAX", 
            "website_slug": "paxos-standard-token", 
            "rank": 52, 
            "circulating_supply": 119149035.0, 
            "total_supply": 119268723.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.0108471471, 
                    "volume_24h": 54734804.2476905, 
                    "market_cap": 120441462.0, 
                    "percent_change_1h": -0.13, 
                    "percent_change_24h": -0.28, 
                    "percent_change_7d": -0.34
                }
            }, 
            "last_updated": 1553335749
        }, 
        "1521": {
            "id": 1521, 
            "name": "Komodo", 
            "symbol": "KMD", 
            "website_slug": "komodo", 
            "rank": 53, 
            "circulating_supply": 112345814.0, 
            "total_supply": 112345814.0, 
            "max_supply": 200000000.0, 
            "quotes": {
                "USD": {
                    "price": 1.0707596887, 
                    "volume_24h": 899852.678885432, 
                    "market_cap": 120295369.0, 
                    "percent_change_1h": 0.25, 
                    "percent_change_24h": 2.14, 
                    "percent_change_7d": -0.36
                }
            }, 
            "last_updated": 1553335742
        }, 
        "693": {
            "id": 693, 
            "name": "Verge", 
            "symbol": "XVG", 
            "website_slug": "verge", 
            "rank": 54, 
            "circulating_supply": 15797108020.0, 
            "total_supply": 15797108020.0, 
            "max_supply": 16555000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.0075447479, 
                    "volume_24h": 4345539.84739065, 
                    "market_cap": 119185197.0, 
                    "percent_change_1h": -0.04, 
                    "percent_change_24h": 2.97, 
                    "percent_change_7d": 4.69
                }
            }, 
            "last_updated": 1553335742
        }, 
        "1866": {
            "id": 1866, 
            "name": "Bytom", 
            "symbol": "BTM", 
            "website_slug": "bytom", 
            "rank": 55, 
            "circulating_supply": 1002499275.0, 
            "total_supply": 1407000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.1169548073, 
                    "volume_24h": 6662668.17436789, 
                    "market_cap": 117247109.0, 
                    "percent_change_1h": -0.32, 
                    "percent_change_24h": 0.44, 
                    "percent_change_7d": 13.91
                }
            }, 
            "last_updated": 1553335743
        }, 
        "3115": {
            "id": 3115, 
            "name": "Maximine Coin", 
            "symbol": "MXM", 
            "website_slug": "maximine-coin", 
            "rank": 56, 
            "circulating_supply": 1649000000.0, 
            "total_supply": 16000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0710298197, 
                    "volume_24h": 13277767.9809559, 
                    "market_cap": 117128173.0, 
                    "percent_change_1h": 1.9, 
                    "percent_change_24h": 20.18, 
                    "percent_change_7d": 101.3
                }
            }, 
            "last_updated": 1553335747
        }, 
        "2603": {
            "id": 2603, 
            "name": "Pundi X", 
            "symbol": "NPXS", 
            "website_slug": "pundi-x", 
            "rank": 57, 
            "circulating_supply": 174450657817.0, 
            "total_supply": 274555193861.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0006239311, 
                    "volume_24h": 1814676.56512182, 
                    "market_cap": 108845196.0, 
                    "percent_change_1h": -0.97, 
                    "percent_change_24h": 0.07, 
                    "percent_change_7d": -5.72
                }
            }, 
            "last_updated": 1553335746
        }, 
        "1042": {
            "id": 1042, 
            "name": "Siacoin", 
            "symbol": "SC", 
            "website_slug": "siacoin", 
            "rank": 58, 
            "circulating_supply": 39955802097.0, 
            "total_supply": 39955802097.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0026834033, 
                    "volume_24h": 1839744.46616224, 
                    "market_cap": 107217533.0, 
                    "percent_change_1h": 0.08, 
                    "percent_change_24h": 0.41, 
                    "percent_change_7d": -3.95
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2405": {
            "id": 2405, 
            "name": "IOST", 
            "symbol": "IOST", 
            "website_slug": "iostoken", 
            "rank": 59, 
            "circulating_supply": 12013965609.0, 
            "total_supply": 21000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0087803986, 
                    "volume_24h": 31957927.973719, 
                    "market_cap": 105487407.0, 
                    "percent_change_1h": -1.13, 
                    "percent_change_24h": -2.76, 
                    "percent_change_7d": 12.73
                }
            }, 
            "last_updated": 1553335746
        }, 
        "2416": {
            "id": 2416, 
            "name": "THETA", 
            "symbol": "THETA", 
            "website_slug": "theta", 
            "rank": 60, 
            "circulating_supply": 870502690.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.1207618084, 
                    "volume_24h": 17009219.979447, 
                    "market_cap": 105123479.0, 
                    "percent_change_1h": 0.02, 
                    "percent_change_24h": -2.99, 
                    "percent_change_7d": -7.08
                }
            }, 
            "last_updated": 1553335745
        }, 
        "2087": {
            "id": 2087, 
            "name": "KuCoin Shares", 
            "symbol": "KCS", 
            "website_slug": "kucoin-shares", 
            "rank": 61, 
            "circulating_supply": 89939916.0, 
            "total_supply": 179939916.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.1587776824, 
                    "volume_24h": 3956131.64264441, 
                    "market_cap": 104220367.0, 
                    "percent_change_1h": -0.33, 
                    "percent_change_24h": 8.0, 
                    "percent_change_7d": 44.58
                }
            }, 
            "last_updated": 1553335743
        }, 
        "3437": {
            "id": 3437, 
            "name": "ABBC Coin", 
            "symbol": "ABBC", 
            "website_slug": "abbc-coin", 
            "rank": 62, 
            "circulating_supply": 457480717.0, 
            "total_supply": 1002164670.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.2071424889, 
                    "volume_24h": 37278937.8514638, 
                    "market_cap": 94763694.0, 
                    "percent_change_1h": -0.19, 
                    "percent_change_24h": -13.84, 
                    "percent_change_7d": -40.75
                }
            }, 
            "last_updated": 1553335748
        }, 
        "2874": {
            "id": 2874, 
            "name": "Aurora", 
            "symbol": "AOA", 
            "website_slug": "aurora", 
            "rank": 63, 
            "circulating_supply": 6542330148.0, 
            "total_supply": 10000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0144365896, 
                    "volume_24h": 2252290.87799623, 
                    "market_cap": 94448935.0, 
                    "percent_change_1h": -0.41, 
                    "percent_change_24h": 8.93, 
                    "percent_change_7d": -3.87
                }
            }, 
            "last_updated": 1553335747
        }, 
        "1343": {
            "id": 1343, 
            "name": "Stratis", 
            "symbol": "STRAT", 
            "website_slug": "stratis", 
            "rank": 64, 
            "circulating_supply": 99259680.0, 
            "total_supply": 99259680.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.9139674933, 
                    "volume_24h": 1167787.21423056, 
                    "market_cap": 90720121.0, 
                    "percent_change_1h": 0.27, 
                    "percent_change_24h": 2.18, 
                    "percent_change_7d": -4.69
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2308": {
            "id": 2308, 
            "name": "Dai", 
            "symbol": "DAI", 
            "website_slug": "dai", 
            "rank": 65, 
            "circulating_supply": 90914172.0, 
            "total_supply": 90914172.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.9943475211, 
                    "volume_24h": 39113768.4145669, 
                    "market_cap": 90400282.0, 
                    "percent_change_1h": 0.3, 
                    "percent_change_24h": 0.16, 
                    "percent_change_7d": -0.37
                }
            }, 
            "last_updated": 1553335745
        }, 
        "2900": {
            "id": 2900, 
            "name": "Project Pai", 
            "symbol": "PAI", 
            "website_slug": "project-pai", 
            "rank": 66, 
            "circulating_supply": 1450682440.0, 
            "total_supply": 1583536500.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0577285973, 
                    "volume_24h": 6230171.50615869, 
                    "market_cap": 83745862.0, 
                    "percent_change_1h": 0.53, 
                    "percent_change_24h": 11.17, 
                    "percent_change_7d": 23.21
                }
            }, 
            "last_updated": 1553335747
        }, 
        "1759": {
            "id": 1759, 
            "name": "Status", 
            "symbol": "SNT", 
            "website_slug": "status", 
            "rank": 67, 
            "circulating_supply": 3470483788.0, 
            "total_supply": 6804870174.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0232031903, 
                    "volume_24h": 10243875.1334149, 
                    "market_cap": 80526296.0, 
                    "percent_change_1h": -0.17, 
                    "percent_change_24h": 3.39, 
                    "percent_change_7d": 0.91
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1789": {
            "id": 1789, 
            "name": "Populous", 
            "symbol": "PPT", 
            "website_slug": "populous", 
            "rank": 68, 
            "circulating_supply": 53252246.0, 
            "total_supply": 53252246.0, 
            "max_supply": 53252246.0, 
            "quotes": {
                "USD": {
                    "price": 1.4443301902, 
                    "volume_24h": 985299.454416643, 
                    "market_cap": 76913827.0, 
                    "percent_change_1h": 0.2, 
                    "percent_change_24h": 2.2, 
                    "percent_change_7d": 2.04
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1455": {
            "id": 1455, 
            "name": "Golem", 
            "symbol": "GNT", 
            "website_slug": "golem-network-tokens", 
            "rank": 69, 
            "circulating_supply": 963622000.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0794976497, 
                    "volume_24h": 1551549.85406224, 
                    "market_cap": 76605684.0, 
                    "percent_change_1h": 0.23, 
                    "percent_change_24h": 1.01, 
                    "percent_change_7d": 4.76
                }
            }, 
            "last_updated": 1553335742
        }, 
        "3116": {
            "id": 3116, 
            "name": "Insight Chain", 
            "symbol": "INB", 
            "website_slug": "insight-chain", 
            "rank": 70, 
            "circulating_supply": 349902689.0, 
            "total_supply": 10000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.2179936507, 
                    "volume_24h": 3444639.78039871, 
                    "market_cap": 76276565.0, 
                    "percent_change_1h": -0.45, 
                    "percent_change_24h": 0.64, 
                    "percent_change_7d": -19.76
                }
            }, 
            "last_updated": 1553335747
        }, 
        "1320": {
            "id": 1320, 
            "name": "Ardor", 
            "symbol": "ARDR", 
            "website_slug": "ardor", 
            "rank": 71, 
            "circulating_supply": 998999495.0, 
            "total_supply": 998999495.0, 
            "max_supply": 998999495.0, 
            "quotes": {
                "USD": {
                    "price": 0.0714008207, 
                    "volume_24h": 878463.266346529, 
                    "market_cap": 71329384.0, 
                    "percent_change_1h": 0.46, 
                    "percent_change_24h": 3.03, 
                    "percent_change_7d": 3.54
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2135": {
            "id": 2135, 
            "name": "Revain", 
            "symbol": "R", 
            "website_slug": "revain", 
            "rank": 72, 
            "circulating_supply": 484450000.0, 
            "total_supply": 484450000.0, 
            "max_supply": 484450000.0, 
            "quotes": {
                "USD": {
                    "price": 0.1427440673, 
                    "volume_24h": 784999.475102919, 
                    "market_cap": 69152363.0, 
                    "percent_change_1h": -0.31, 
                    "percent_change_24h": 0.74, 
                    "percent_change_7d": 3.51
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1586": {
            "id": 1586, 
            "name": "Ark", 
            "symbol": "ARK", 
            "website_slug": "ark", 
            "rank": 73, 
            "circulating_supply": 109140366.0, 
            "total_supply": 140390366.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.6315431285, 
                    "volume_24h": 645060.696213318, 
                    "market_cap": 68926848.0, 
                    "percent_change_1h": -0.2, 
                    "percent_change_24h": 0.5, 
                    "percent_change_7d": -3.72
                }
            }, 
            "last_updated": 1553335742
        }, 
        "1750": {
            "id": 1750, 
            "name": "GXChain", 
            "symbol": "GXC", 
            "website_slug": "gxchain", 
            "rank": 74, 
            "circulating_supply": 60000000.0, 
            "total_supply": 99627196.0, 
            "max_supply": 100000000.0, 
            "quotes": {
                "USD": {
                    "price": 1.1269820825, 
                    "volume_24h": 8777661.81488737, 
                    "market_cap": 67618925.0, 
                    "percent_change_1h": -0.2, 
                    "percent_change_24h": -0.32, 
                    "percent_change_7d": 23.41
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2829": {
            "id": 2829, 
            "name": "REPO", 
            "symbol": "REPO", 
            "website_slug": "repo", 
            "rank": 75, 
            "circulating_supply": 109958607.0, 
            "total_supply": 356999900.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.5869835226, 
                    "volume_24h": 57971.5218899815, 
                    "market_cap": 64543890.0, 
                    "percent_change_1h": -5.22, 
                    "percent_change_24h": 2.95, 
                    "percent_change_7d": -12.74
                }
            }, 
            "last_updated": 1553335746
        }, 
        "2027": {
            "id": 2027, 
            "name": "Cryptonex", 
            "symbol": "CNX", 
            "website_slug": "cryptonex", 
            "rank": 76, 
            "circulating_supply": 55686329.0, 
            "total_supply": 107135054.0, 
            "max_supply": 210000000.0, 
            "quotes": {
                "USD": {
                    "price": 1.141067684, 
                    "volume_24h": 8926779.62844466, 
                    "market_cap": 63541870.0, 
                    "percent_change_1h": -1.15, 
                    "percent_change_24h": -0.27, 
                    "percent_change_7d": -1.35
                }
            }, 
            "last_updated": 1553335743
        }, 
        "2349": {
            "id": 2349, 
            "name": "Mixin", 
            "symbol": "XIN", 
            "website_slug": "mixin", 
            "rank": 77, 
            "circulating_supply": 434419.0, 
            "total_supply": 1000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 145.641472805, 
                    "volume_24h": 587263.135904176, 
                    "market_cap": 63269427.0, 
                    "percent_change_1h": -0.03, 
                    "percent_change_24h": 1.58, 
                    "percent_change_7d": -3.76
                }
            }, 
            "last_updated": 1553335745
        }, 
        "3306": {
            "id": 3306, 
            "name": "Gemini Dollar", 
            "symbol": "GUSD", 
            "website_slug": "gemini-dollar", 
            "rank": 78, 
            "circulating_supply": 61621381.0, 
            "total_supply": 61621381.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 1.0108153734, 
                    "volume_24h": 21716801.4117324, 
                    "market_cap": 62287839.0, 
                    "percent_change_1h": -0.16, 
                    "percent_change_24h": -0.28, 
                    "percent_change_7d": -0.91
                }
            }, 
            "last_updated": 1553335748
        }, 
        "1903": {
            "id": 1903, 
            "name": "HyperCash", 
            "symbol": "HC", 
            "website_slug": "hypercash", 
            "rank": 79, 
            "circulating_supply": 43529781.0, 
            "total_supply": 43529781.0, 
            "max_supply": 84000000.0, 
            "quotes": {
                "USD": {
                    "price": 1.3705250409, 
                    "volume_24h": 1718174.26733937, 
                    "market_cap": 59658655.0, 
                    "percent_change_1h": 1.28, 
                    "percent_change_24h": 4.62, 
                    "percent_change_7d": 1.67
                }
            }, 
            "last_updated": 1553335743
        }, 
        "291": {
            "id": 291, 
            "name": "MaidSafeCoin", 
            "symbol": "MAID", 
            "website_slug": "maidsafecoin", 
            "rank": 80, 
            "circulating_supply": 452552412.0, 
            "total_supply": 452552412.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.1298652435, 
                    "volume_24h": 261240.255777524, 
                    "market_cap": 58770829.0, 
                    "percent_change_1h": 0.03, 
                    "percent_change_24h": 0.74, 
                    "percent_change_7d": 2.8
                }
            }, 
            "last_updated": 1553335741
        }, 
        "2300": {
            "id": 2300, 
            "name": "WAX", 
            "symbol": "WAX", 
            "website_slug": "wax", 
            "rank": 81, 
            "circulating_supply": 942694871.0, 
            "total_supply": 1850000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0614514724, 
                    "volume_24h": 322506.919241758, 
                    "market_cap": 57929988.0, 
                    "percent_change_1h": 0.01, 
                    "percent_change_24h": 0.88, 
                    "percent_change_7d": -1.05
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2137": {
            "id": 2137, 
            "name": "Electroneum", 
            "symbol": "ETN", 
            "website_slug": "electroneum", 
            "rank": 82, 
            "circulating_supply": 9149407051.0, 
            "total_supply": 9149407051.0, 
            "max_supply": 21000000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.0062898633, 
                    "volume_24h": 9351252.18755502, 
                    "market_cap": 57548520.0, 
                    "percent_change_1h": 0.01, 
                    "percent_change_24h": -2.14, 
                    "percent_change_7d": -4.09
                }
            }, 
            "last_updated": 1553335744
        }, 
        "1087": {
            "id": 1087, 
            "name": "Factom", 
            "symbol": "FCT", 
            "website_slug": "factom", 
            "rank": 83, 
            "circulating_supply": 9410899.0, 
            "total_supply": 9410899.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 6.0650791873, 
                    "volume_24h": 135862.120978432, 
                    "market_cap": 57077845.0, 
                    "percent_change_1h": -0.46, 
                    "percent_change_24h": -2.77, 
                    "percent_change_7d": -9.35
                }
            }, 
            "last_updated": 1553335741
        }, 
        "2588": {
            "id": 2588, 
            "name": "Loom Network", 
            "symbol": "LOOM", 
            "website_slug": "loom-network", 
            "rank": 84, 
            "circulating_supply": 763436089.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0734302988, 
                    "volume_24h": 3142725.40120225, 
                    "market_cap": 56059340.0, 
                    "percent_change_1h": 0.52, 
                    "percent_change_24h": -2.07, 
                    "percent_change_7d": 8.9
                }
            }, 
            "last_updated": 1553335746
        }, 
        "2213": {
            "id": 2213, 
            "name": "QASH", 
            "symbol": "QASH", 
            "website_slug": "qash", 
            "rank": 85, 
            "circulating_supply": 350000000.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.1575308812, 
                    "volume_24h": 135404.748982484, 
                    "market_cap": 55135808.0, 
                    "percent_change_1h": -0.39, 
                    "percent_change_24h": 1.93, 
                    "percent_change_7d": 3.38
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2772": {
            "id": 2772, 
            "name": "Digitex Futures", 
            "symbol": "DGTX", 
            "website_slug": "digitex-futures", 
            "rank": 86, 
            "circulating_supply": 737500000.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0728828005, 
                    "volume_24h": 1026075.88133525, 
                    "market_cap": 53751065.0, 
                    "percent_change_1h": -0.1, 
                    "percent_change_24h": -2.7, 
                    "percent_change_7d": -15.52
                }
            }, 
            "last_updated": 1553335746
        }, 
        "1966": {
            "id": 1966, 
            "name": "Decentraland", 
            "symbol": "MANA", 
            "website_slug": "decentraland", 
            "rank": 87, 
            "circulating_supply": 1050141509.0, 
            "total_supply": 2644403343.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0505756129, 
                    "volume_24h": 2888115.31860478, 
                    "market_cap": 53111550.0, 
                    "percent_change_1h": 0.05, 
                    "percent_change_24h": 1.3, 
                    "percent_change_7d": 2.97
                }
            }, 
            "last_updated": 1553335744
        }, 
        "1934": {
            "id": 1934, 
            "name": "Loopring", 
            "symbol": "LRC", 
            "website_slug": "loopring", 
            "rank": 88, 
            "circulating_supply": 828954240.0, 
            "total_supply": 1374955752.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0626357393, 
                    "volume_24h": 1061978.57759421, 
                    "market_cap": 51922162.0, 
                    "percent_change_1h": -0.04, 
                    "percent_change_24h": 2.09, 
                    "percent_change_7d": -3.66
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1776": {
            "id": 1776, 
            "name": "Crypto.com", 
            "symbol": "MCO", 
            "website_slug": "crypto-com", 
            "rank": 89, 
            "circulating_supply": 15793831.0, 
            "total_supply": 31587682.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 3.2860599257, 
                    "volume_24h": 2691954.12762571, 
                    "market_cap": 51899475.0, 
                    "percent_change_1h": -0.04, 
                    "percent_change_24h": 1.52, 
                    "percent_change_7d": -0.13
                }
            }, 
            "last_updated": 1553335743
        }, 
        "1925": {
            "id": 1925, 
            "name": "Waltonchain", 
            "symbol": "WTC", 
            "website_slug": "waltonchain", 
            "rank": 90, 
            "circulating_supply": 41007759.0, 
            "total_supply": 70000000.0, 
            "max_supply": 100000000.0, 
            "quotes": {
                "USD": {
                    "price": 1.2644663573, 
                    "volume_24h": 3313676.55856806, 
                    "market_cap": 51852932.0, 
                    "percent_change_1h": 0.31, 
                    "percent_change_24h": 0.09, 
                    "percent_change_7d": -6.02
                }
            }, 
            "last_updated": 1553335743
        }, 
        "3224": {
            "id": 3224, 
            "name": "Qubitica", 
            "symbol": "QBIT", 
            "website_slug": "qubitica", 
            "rank": 91, 
            "circulating_supply": 2805292.0, 
            "total_supply": 10000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 17.7774581085, 
                    "volume_24h": 84003.4758558728, 
                    "market_cap": 49870967.0, 
                    "percent_change_1h": -0.13, 
                    "percent_change_24h": -0.36, 
                    "percent_change_7d": -2.16
                }
            }, 
            "last_updated": 1553335748
        }, 
        "1414": {
            "id": 1414, 
            "name": "Zcoin", 
            "symbol": "XZC", 
            "website_slug": "zcoin", 
            "rank": 92, 
            "circulating_supply": 7049883.0, 
            "total_supply": 21400000.0, 
            "max_supply": 21400000.0, 
            "quotes": {
                "USD": {
                    "price": 6.9330693009, 
                    "volume_24h": 2329069.27924742, 
                    "market_cap": 48877330.0, 
                    "percent_change_1h": 0.59, 
                    "percent_change_24h": 4.19, 
                    "percent_change_7d": 2.45
                }
            }, 
            "last_updated": 1553335742
        }, 
        "2299": {
            "id": 2299, 
            "name": "aelf", 
            "symbol": "ELF", 
            "website_slug": "aelf", 
            "rank": 93, 
            "circulating_supply": 280000000.0, 
            "total_supply": 300000000.0, 
            "max_supply": 1000000000.0, 
            "quotes": {
                "USD": {
                    "price": 0.1742485943, 
                    "volume_24h": 4334832.6413079, 
                    "market_cap": 48789606.0, 
                    "percent_change_1h": -0.25, 
                    "percent_change_24h": 0.91, 
                    "percent_change_7d": 2.27
                }
            }, 
            "last_updated": 1553335745
        }, 
        "1169": {
            "id": 1169, 
            "name": "PIVX", 
            "symbol": "PIVX", 
            "website_slug": "pivx", 
            "rank": 94, 
            "circulating_supply": 56781166.0, 
            "total_supply": 56781166.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.8572490701, 
                    "volume_24h": 642155.372444222, 
                    "market_cap": 48675602.0, 
                    "percent_change_1h": 2.41, 
                    "percent_change_24h": 5.75, 
                    "percent_change_7d": 1.45
                }
            }, 
            "last_updated": 1553335742
        }, 
        "3144": {
            "id": 3144, 
            "name": "ThoreCoin", 
            "symbol": "THR", 
            "website_slug": "thorecoin", 
            "rank": 95, 
            "circulating_supply": 86686.0, 
            "total_supply": 100000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 557.318499432, 
                    "volume_24h": 83238.4421123775, 
                    "market_cap": 48311723.0, 
                    "percent_change_1h": 0.04, 
                    "percent_change_24h": 0.22, 
                    "percent_change_7d": 1.21
                }
            }, 
            "last_updated": 1553335748
        }, 
        "1982": {
            "id": 1982, 
            "name": "Kyber Network", 
            "symbol": "KNC", 
            "website_slug": "kyber-network", 
            "rank": 96, 
            "circulating_supply": 164433012.0, 
            "total_supply": 215014561.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.2710064309, 
                    "volume_24h": 6361186.23661844, 
                    "market_cap": 44562404.0, 
                    "percent_change_1h": -0.14, 
                    "percent_change_24h": 2.13, 
                    "percent_change_7d": 10.18
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2694": {
            "id": 2694, 
            "name": "Nexo", 
            "symbol": "NEXO", 
            "website_slug": "nexo", 
            "rank": 97, 
            "circulating_supply": 560000011.0, 
            "total_supply": 1000000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.0786247961, 
                    "volume_24h": 7786846.85868852, 
                    "market_cap": 44029887.0, 
                    "percent_change_1h": 0.21, 
                    "percent_change_24h": -2.7, 
                    "percent_change_7d": -7.75
                }
            }, 
            "last_updated": 1553335746
        }, 
        "2606": {
            "id": 2606, 
            "name": "Wanchain", 
            "symbol": "WAN", 
            "website_slug": "wanchain", 
            "rank": 98, 
            "circulating_supply": 106152493.0, 
            "total_supply": 210000000.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.4138517547, 
                    "volume_24h": 3099029.59836351, 
                    "market_cap": 43931395.0, 
                    "percent_change_1h": -0.64, 
                    "percent_change_24h": 6.56, 
                    "percent_change_7d": -5.27
                }
            }, 
            "last_updated": 1553335746
        }, 
        "2092": {
            "id": 2092, 
            "name": "NULS", 
            "symbol": "NULS", 
            "website_slug": "nuls", 
            "rank": 99, 
            "circulating_supply": 64562561.0, 
            "total_supply": 103485843.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.677886419, 
                    "volume_24h": 22289197.7508697, 
                    "market_cap": 43766083.0, 
                    "percent_change_1h": -0.38, 
                    "percent_change_24h": 0.69, 
                    "percent_change_7d": 17.2
                }
            }, 
            "last_updated": 1553335744
        }, 
        "2403": {
            "id": 2403, 
            "name": "MOAC", 
            "symbol": "MOAC", 
            "website_slug": "moac", 
            "rank": 100, 
            "circulating_supply": 62463334.0, 
            "total_supply": 151205864.0, 
            "max_supply": null, 
            "quotes": {
                "USD": {
                    "price": 0.6989068988, 
                    "volume_24h": 52773.6419830923, 
                    "market_cap": 43656055.0, 
                    "percent_change_1h": 0.28, 
                    "percent_change_24h": 2.8, 
                    "percent_change_7d": -8.36
                }
            }, 
            "last_updated": 1553335745
        }
    }, 
    "metadata": {
        "timestamp": 1553334618, 
        "warning": "WARNING: This API is now deprecated and will be taken offline soon.  Please switch to the new CoinMarketCap API to avoid interruptions in service. (https://pro.coinmarketcap.com/migrate/)", 
        "num_cryptocurrencies": 2121, 
        "error": null
    }
}`)
)

type(
	TopRank struct{
		Id int64	`sql:"primary;name:id"			json:"id"		bson:"id"`
		Type int8	`sql:"primary;name:type"		json:"type"		bson:"type"`
		Name string `sql:"name:name"				json:"name"		bson:"name"`
		Score int `sql:"name:score"					json:"score"	bson:"score"`
		Value [2]int `sql:"name:value"				json:"value"	bson:"value"`
		LastTime int64 `sql:"datetime;name:last_time"	json:"last_time"	bson:"last_time"`
	}
)

var(
	ntimes = 10000
)

func TestJson(t *testing.T){
	for i := 0; i < ntimes; i++{
		json.Marshal(&Coinmarketcap{})
	}
}

func TestUJson(t *testing.T){
	for i := 0; i < ntimes; i++{
		json.Unmarshal(data, &Coinmarketcap{})
	}
}

func TestJJson(t *testing.T){
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < ntimes; i++{
		json.Marshal(&Coinmarketcap{})
	}
}

func TestUJJson(t *testing.T){
	aa := &Coinmarketcap{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < ntimes; i++{
		json.Unmarshal(data, &aa)
	}
	data1, _ = bson.Marshal(aa)
}

func TestBson(t *testing.T){
	for i := 0; i < ntimes; i++{
		bson.Marshal(&Coinmarketcap{})
	}
}

func TestUbson(t *testing.T){
	aa := &Coinmarketcap{}
	for i := 0; i < ntimes; i++{
		bson.Unmarshal(data1, &aa)
	}
}