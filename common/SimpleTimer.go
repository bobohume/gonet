package common

import (
	"time"
)

type (
	SimpleTimer struct {
		interval int
		lastTime int64
		count    int
		isActive bool
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

func (s *SimpleTimer) SetInterval(interval int) {
	s.interval = interval
}

func (s *SimpleTimer) CheckTimer() bool {
	if !s.isActive {
		return false
	}

	curTime := time.Now().Unix()
	if curTime-s.lastTime >= int64(s.interval) {
		s.lastTime = curTime
		s.count++
		return true
	}
	return false
}

func (s *SimpleTimer) Start() {
	s.lastTime = time.Now().Unix()
	s.count = 0
	s.isActive = true
}

func (s *SimpleTimer) Stop() {
	s.lastTime = time.Now().Unix()
	s.count = 0
	s.isActive = false
}

func (s *SimpleTimer) GetTimerCount() int {
	return s.count
}

func (s *SimpleTimer) IsActived() bool {
	return s.isActive
}

func NewSimpleTimer(interval int) *SimpleTimer {
	simpleTimer := new(SimpleTimer)
	simpleTimer.SetInterval(interval)
	return simpleTimer
}
