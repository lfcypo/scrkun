package notifiy

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lfcypo/scrkun/logger"
	"github.com/lfcypo/scrkun/usagedt"
)

var log = logger.New("Notify")

var Disable = fmt.Errorf("notifier disabled")

type Notifier interface {
	Notify(msg string) error
}

var notifiers = []Notifier{
	&Bark{},
	&Xtuis{},
}

func Notify(dts []*usagedt.Usage, machineName string) {
	msg := GenerateMessage(dts, machineName)
	log.Infof("推送消息: %s", msg)
	for _, notifier := range notifiers {
		err := notifier.Notify(msg)
		if err != nil {
			if errors.Is(err, Disable) {
				continue
			}
			log.Warnf("发送消息推送失败 请检查网络与配置文件是否正确: %v", err)
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
