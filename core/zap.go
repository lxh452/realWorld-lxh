package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"realWorld/global"
)

type ZapCore struct {
	level zapcore.Level
	zapcore.Core
}

func NewZapCore(level zapcore.Level) *ZapCore {
	entity := &ZapCore{level: level}
	syncer := entity.WriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	})
	entity.Core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), syncer, levelEnabler)
	return entity
}

func (z *ZapCore) WriteSyncer() zapcore.WriteSyncer {
	cutter := NewCutter(
		CutterWithLayout(global.CONFIG.Logs.Layout),
		CutterWithLevel(z.level),
		CutterWithDirector(global.CONFIG.Logs.Dir),
	)
	return zapcore.AddSync(cutter)
}
