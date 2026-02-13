package model

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model `json:"-"`

	// Time 事件时间
	Time time.Time `json:"time"`
	// EventType 事件类型
	EventType EventType `json:"eventType"`
	// Reason 判断原因
	Reason string `json:"reason"`
	// UsageType 用法类型
	UsageType UsageType `json:"usageType"`
	// Weight 权重
	Weight float64 `json:"weight"`

	// ExtraInfo 附加信息
	ExtraInfo string `json:"extraInfo"`
}

func NewEvent(usageType UsageType, reason string, weight float64, extraInfo string) *Event {
	timestamp := time.Now()

	return &Event{
		Time:      timestamp,
		UsageType: usageType,
		Reason:    reason,
		ExtraInfo: extraInfo,
		Weight:    weight,
	}
}

func (e *Event) NormalEvent() {
	e.EventType = NormalEvent
}

func (e *Event) AbnormalEvent() {
	e.EventType = AbnormalEvent
}

func (e *Event) IsNormalEvent() bool {
	return e.EventType == NormalEvent
}

func (e *Event) IsAbnormalEvent() bool {
	return !e.IsNormalEvent()
}

func (e *Event) Notified() bool {
	return e.IsAbnormalEvent()
}
