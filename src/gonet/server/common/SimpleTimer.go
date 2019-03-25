package common

import (
	"time"
)

type(
	SimpleTimer struct{
		m_interval 	int
		m_lastTime 	int64
		m_count 	int
		m_isActive 	bool
	}

	ISimpleTimer interface {
		SetInterval(int)
		CheckTimer() bool
		Start()
		Stop()
		GetTimerCount() int
		IsActived() bool
	}
)

func (this *SimpleTimer) SetInterval(interval int){
	this.m_interval = interval
}

func (this *SimpleTimer)  CheckTimer() bool {
	if !this.m_isActive{
		return false
	}

	curTime := time.Now().Unix()
	if curTime - this.m_lastTime >= int64(this.m_interval){
		this.m_lastTime = curTime
		this.m_count++
		return true
	}
	return false
}

func (this *SimpleTimer)  Start() {
	this.m_lastTime = time.Now().Unix()
	this.m_count    = 0;
	this.m_isActive = true;
}

func (this *SimpleTimer) Stop() {
	this.m_lastTime = time.Now().Unix()
	this.m_count    = 0;
	this.m_isActive = false;
}

func (this *SimpleTimer) GetTimerCount() int {
	return this.m_count
}

func (this *SimpleTimer) IsActived() bool {
	return this.m_isActive
}

func NewSimpleTimer(interval int) *SimpleTimer{
	simpleTimer := new(SimpleTimer)
	simpleTimer.SetInterval(interval)
	return simpleTimer
}

