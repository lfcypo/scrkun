package notifiy

import (
	"fmt"
	"strings"
	"time"

	"github.com/lfcypo/scrkun/logger"
	"github.com/lfcypo/scrkun/usagedt"
)

var log = logger.New("Notify")

type Notifier interface {
	Notify(msg string) error
}

var notifiers = []Notifier{
	&Bark{},
}

func Notify(dts []*usagedt.Usage, machineName string) {
	msg := GenerateMessage(dts, machineName)
	for _, notifier := range notifiers {
		err := notifier.Notify(msg)
		if err != nil {
			log.Warnf("Failed to notify: %v", err)
		}
	}
}

func GenerateMessage(dts []*usagedt.Usage, machineName string) string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("于 %s 在设备 %s 上监测到异常行为: ", timeStr, machineName))
	for _, dt := range dts {
		s := fmt.Sprintf("%s", dt.Reason)
		sb.WriteString(s)
	}

	sb.WriteString(" 请处理 谢谢")

	return sb.String()
}
