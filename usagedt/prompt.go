package usagedt

import (
	_ "embed"
	"sync"
)

//go:embed sysprompt.txt
var rawSysPrompt string

var sysPrompt string

var onceDoGenerateSysPrompt = &sync.Once{}

func SysPrompt() string {
	onceDoGenerateSysPrompt.Do(func() {
		sysPrompt = rawSysPrompt
	})
	return sysPrompt
}
