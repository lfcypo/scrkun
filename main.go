package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lfcypo/scrkun/capture"
	_ "github.com/lfcypo/scrkun/config"
	"github.com/lfcypo/scrkun/logger"
	"github.com/lfcypo/scrkun/notifiy"
	"github.com/lfcypo/scrkun/usagedt"
	"github.com/lfcypo/scrkun/util"
	"github.com/lfcypo/viperx"
	"github.com/spf13/viper"
)

var log = logger.New("Main")

var shouldNotifyUsage = []string{
	"video",
	"game",
	"novel",
}

var contains = util.ContainsFunc(shouldNotifyUsage)

//func init() {
//	logger.DisableColor()
//}

func main() {

	sleep := time.Duration(viperx.GetInt("monitor.sleep", 5))
	delay := time.Second * time.Duration(viperx.GetInt("monitor.interval", 300))

	lastDetectedAt := time.UnixMicro(0)
	i := 0
	for {
		i++

		time.Sleep(time.Second * sleep)
		timeDiff := time.Since(lastDetectedAt)
		if timeDiff < delay {
			log.Infof("Within the sleep time, there are %.2f seconds before the next detection", delay.Seconds()-timeDiff.Seconds())
			continue
		}

		img, err := capture.Capture(fmt.Sprintf("tmp/%d-%s.png", i, time.Now().Format("2006-01-02-15-04-05")))
		if err != nil {
			if errors.Is(err, capture.Skip) {
				continue
			}
			log.Errorf("Capture error: %v", err)
		}
		log.Infof("Capture success: %s", img.Bounds())

		usages := usagedt.Detect(img)
		if usages == nil {
			continue
		}
		log.Infof("Usage: %s", usages.String())

		var possibleUsage []*usagedt.Usage
		for _, usage := range usages.Usages {
			if usage.Weight < viperx.GetFloat64("detect.threshold", 0.5) || !contains(usage.Usage) {
				continue
			}
			lastDetectedAt = time.Now()
			possibleUsage = append(possibleUsage, usage)
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
