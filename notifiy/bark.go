package notifiy

import (
	"errors"
	"sync"

	"github.com/jzksnsjswkw/go-bark"
	"github.com/spf13/viper"
)

var doOnceInitBark = &sync.Once{}

var token string
var enable = false

func init() {
	doOnceInitBark.Do(func() {
		if !viper.IsSet("bark.token") {
			return
		}
		token = viper.GetString("bark.token")
		enable = true
	})
}

type Bark struct {
}

func (b Bark) Notify(msg string) error {
	if !enable {
		return errors.New("bark notifier is disable by config")
	}

	err := bark.Push(&bark.Options{
		Msg:   msg,
		Token: token,
	})
	return err
}
