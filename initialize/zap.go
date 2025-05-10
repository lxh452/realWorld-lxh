package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ccore "realWorld/core"
	"realWorld/global"
)

func MustLoadZap() {
	levels := Levels()
	fmt.Println("等级", levels)
	length := len(levels)
	fmt.Print("levels:", levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := ccore.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger := zap.New(zapcore.NewTee(cores...))
	global.Logger = logger
}

// Levels 根据字符串转化为 zapcore.Levels
func Levels() []zapcore.Level {
	levels := make([]zapcore.Level, 0, 7)
	level, err := zapcore.ParseLevel("info")
	if err != nil {
		level = zapcore.DebugLevel
	}
	for ; level <= zapcore.FatalLevel; level++ {
		levels = append(levels, level)
	}
	fmt.Println(levels)
	return levels
}
