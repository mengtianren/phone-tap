package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CloseOffsetX     int  `yaml:"close_offsetx"`
	CloseOffsetY     int  `yaml:"close_offsety"`
	SuccessOffsetX   int  `yaml:"success_offsetx"`
	SuccessOffsetY   int  `yaml:"success_offsety"`
	AwaitTime        int  `yaml:"await_time"`
	EndTime          int  `yaml:"end_time"`
	LiveCloseTime    int  `yaml:"live_close_time"`
	LiveCloseOffsetX int  `yaml:"live_close_offsetx"`
	LiveCloseOffsetY int  `yaml:"live_close_offsety"`
	LiveCloseStart   bool `yaml:"live_close_start"`
}

var Cfg Config

func loadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Print("--------❌ 读取变量失败1，使用默认配置--------\n")
		Cfg = Config{
			CloseOffsetX:   400,
			CloseOffsetY:   0,
			SuccessOffsetX: 0,
			SuccessOffsetY: 0,
			AwaitTime:      3,
		}
		return
	}

	err1 := yaml.Unmarshal(data, &Cfg)
	if err1 != nil {
		fmt.Print("--------❌ 读取变量失败2: ，使用默认配置--------\n")
		Cfg = Config{
			CloseOffsetX:   400,
			CloseOffsetY:   0,
			SuccessOffsetX: 0,
			SuccessOffsetY: 0,
			AwaitTime:      3,
		}
		return
	}
	fmt.Print("--------✅ 读取变量成功，正在启动中-----------\n")

}

func init() {
	loadConfig("./config.yaml")
}
