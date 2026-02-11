package notifiy

import (
	"sync"

	"github.com/jzksnsjswkw/go-bark"
	"github.com/spf13/viper"
)

var doOnceInitBark = &sync.Once{}

var barkToken string
var barkEnable = false

func init() {
	doOnceInitBark.Do(func() {
		if !viper.IsSet("bark.token") {
			return
		}
		barkToken = viper.GetString("bark.token")
		barkEnable = true
	})
}

type Bark struct {
}

func (b Bark) Notify(msg string) error {
	if !barkEnable {
		return Disable
	}

	err := bark.Push(&bark.Options{
		Msg:   msg,
		Token: barkToken,
	})
	return err
}
