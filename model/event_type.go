package model

type EventType = int

const (
	// NormalEvent 正常事件
	NormalEvent EventType = iota
	// AbnormalEvent 异常事件
	AbnormalEvent
)
