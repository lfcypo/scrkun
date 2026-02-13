package model

type UsageType int

// [study: 校本课程学习, code: 写代码或进行编程任务, video: 观看视频, game: 玩游戏或查看相关内容, document: 文件编辑, web: 浏览网页, novel: 阅读小说, empty: 闲置]

const (
	Study UsageType = iota
	Code
	Video
	Game
	Document
	Web
	Novel
	Empty
)

var convMapFromUsageType = map[UsageType]string{
	Study:    "study",
	Code:     "code",
	Video:    "video",
	Game:     "game",
	Document: "document",
	Web:      "web",
	Novel:    "novel",
	Empty:    "empty",
}

func ConvFromUsageType(actionType UsageType) string {
	return convMapFromUsageType[actionType]
}

func (u UsageType) String() string {
	return convMapFromUsageType[u]
}

var convMapToUsageType = map[string]UsageType{
	"study":    Study,
	"code":     Code,
	"video":    Video,
	"game":     Game,
	"document": Document,
	"web":      Web,
	"novel":    Novel,
	"empty":    Empty,
}

func ConvToUsageType(actionType string) UsageType {
	return convMapToUsageType[actionType]
}
