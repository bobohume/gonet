package base

import (
	"time"
)

const(
	INTERVAL_DAY = iota		//当前时间的第二天
	INTERVAL_WEEK = iota		//当前时间的下周一
	INTERVAL_MONTH = iota	//当前时间的下一月第一天)
	INTERVAL_YEAR   = iota   //当前时间的下一年第一天)
	TIME_SET_MAX_VAL = iota //类型最大值
)

func GetNextTime(intervalType int) time.Time{
	t := time.Now()
	if intervalType == INTERVAL_YEAR{
		t = t.AddDate(1, 0, 0)
	}else if intervalType == INTERVAL_MONTH{
		t = t.AddDate(0, 1, 0)
	}else if intervalType == INTERVAL_WEEK{
		if t.Weekday() != time.Sunday{
			t = t.AddDate(0, 0, (8-int(t.Weekday())))
		}else{
			t = t.AddDate(0,0, 1)
		}
	}else{
		t = t.AddDate(0,0, 1)
	}

	DefaultTimeLoc := time.Local
	return  time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, DefaultTimeLoc)
}
