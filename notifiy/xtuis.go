package notifiy

import (
	"sync"

	"github.com/lfcypo/go-xtuis"
	"github.com/spf13/viper"
)

var doOnceInitXtuis = &sync.Once{}

var xtuisToken string
var xtuisEnable = false

func init() {
	doOnceInitXtuis.Do(func() {
		if !viper.IsSet("xtuis.token") {
			return
		}
		xtuisToken = viper.GetString("xtuis.token")
		xtuisEnable = true
	})
}

type Xtuis struct {
}

func (x Xtuis) Notify(msg string) error {
	if !xtuisEnable {
		return Disable
	}

	client := xtuis.NewClient(xtuisToken)
	payload := xtuis.NewPayload("电教平台 AI 智能监控系统 报警", msg)
	err := client.Send(payload)

	return err
}
