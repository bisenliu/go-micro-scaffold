package timezone

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/fx"

	"common/config"
)

var (
	once sync.Once
)

// Params 时区依赖参数
type Params struct {
	fx.In
	Config *config.Config
}

// Init 初始化全局时区设置，只执行一次
func InitTimezone(params Params) {
	once.Do(func() {
		// 从配置中获取时区，如果不存在则默认为 "Asia/Shanghai"
		timezone := "Asia/Shanghai"
		if params.Config.System.Timezone != "" {
			timezone = params.Config.System.Timezone
		}

		// 加载时区
		tz, err := time.LoadLocation(timezone)
		if err != nil {
			panic(fmt.Errorf("load location error: %w", err))
		}

		// 全局设置时区
		time.Local = tz
		os.Setenv("TZ", timezone)
	})
}

// Module FX模块
var Module = fx.Module("timezone",
	fx.Invoke(InitTimezone),
)
