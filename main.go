package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lfcypo/scrkun/capture"
	_ "github.com/lfcypo/scrkun/config"
	"github.com/lfcypo/scrkun/logger"
	"github.com/lfcypo/scrkun/model"
	"github.com/lfcypo/scrkun/notifiy"
	"github.com/lfcypo/scrkun/server"
	"github.com/lfcypo/scrkun/store"
	"github.com/lfcypo/scrkun/usagedt"
	"github.com/lfcypo/scrkun/util"
	"github.com/lfcypo/viperx"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var log = logger.New("Main")
var storeLog = logger.New("Store")

var shouldNotifyUsage = []string{
	"video",
	"game",
	"novel",
}

var db = store.GetDatabase()
var abnormal = util.ContainsFunc(shouldNotifyUsage)

func main() {

	logger.DisableColor()

	go server.StartWebServer(viperx.GetInt("server.port", 8890))
	if util.Contains(os.Args, "-q") {
		select {}
	}

	sleep := time.Duration(viperx.GetInt("monitor.sleep", 5))
	delay := time.Second * time.Duration(viperx.GetInt("monitor.interval", 300))

	lastDetectedAt := time.UnixMicro(0)
	i := 0
	for {
		i++

		time.Sleep(time.Second * sleep)
		timeDiff := time.Since(lastDetectedAt)
		if timeDiff < delay {
			log.Infof("已发出报警 距离下一次检测还有 %.2f 秒", delay.Seconds()-timeDiff.Seconds())
			continue
		}

		img, err := capture.Capture(fmt.Sprintf("tmp/%d-%s.png", i, time.Now().Format("2006-01-02-15-04-05")))
		if err != nil {
			if errors.Is(err, capture.Skip) {
				continue
			}
			log.Errorf("屏幕截图失败: %v", err)
		}
		log.Infof("屏幕截图成功: %s", img.Bounds())

		usages := usagedt.Detect(img)
		if usages == nil {
			continue
		}
		//log.Debugf("检测的用途: %s", usages.String())

		var possibleUsage []*usagedt.Usage
		for _, usage := range usages.Usages {
			if usage.Weight < viperx.GetFloat64("detect.threshold", 0.5) {
				continue
			}

			event := model.NewEvent(model.ConvToUsageType(usage.Usage), usage.Reason, usage.Weight, "")
			log.Infof("检测到使用情况: %s, 权重, %.2f, 原因: %s", usage.Usage, usage.Weight, usage.Reason)
			if abnormal(usage.Usage) {
				log.Infof("需要报告使用情况: %s, 权重, %.2f, 原因: %s", usage.Usage, usage.Weight, usage.Reason)
				lastDetectedAt = time.Now()
				possibleUsage = append(possibleUsage, usage)
				event.AbnormalEvent()
			} else {
				event.NormalEvent()
			}
			err := gorm.G[model.Event](db).Create(context.Background(), event)
			if err != nil {
				storeLog.Errorf("保存事件失败: %v", err)
			}
			storeLog.Info("保存事件成功")
		}

		if len(possibleUsage) > 0 {
			notifiy.Notify(possibleUsage, GetMachineName())
		}
	}

}

func GetMachineName() string {
	if !viper.IsSet("machine.name") {
		hostname, err := os.Hostname()
		if err != nil {
			return "unknown"
		}
		return hostname
	}
	return viper.GetString("machine.name")
}
