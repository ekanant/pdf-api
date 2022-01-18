package logs

import (
	"go.uber.org/zap"
)

var _log *zap.Logger

func init() {
	//logBuilder := zap.NewProductionConfig()
	logBuilder := zap.NewDevelopmentConfig()
	_log, _ = logBuilder.Build(zap.AddCallerSkip(1))
}

func Info(msg interface{}, fields ...zap.Field) {
	switch v := msg.(type) {
	case string:
		{
			_log.Info(v, fields...)
		}
	case error:
		{
			_log.Info(v.Error(), fields...)
		}
	}
}

func Debug(msg interface{}, fields ...zap.Field) {
	switch v := msg.(type) {
	case string:
		{
			_log.Debug(v, fields...)
		}
	case error:
		{
			_log.Debug(v.Error(), fields...)
		}
	}
}

func Error(msg interface{}, fields ...zap.Field) {
	switch v := msg.(type) {
	case string:
		{
			_log.Error(v, fields...)
		}
	case error:
		{
			_log.Error(v.Error(), fields...)
		}
	}
}

func Panic(msg interface{}, fields ...zap.Field) {
	switch v := msg.(type) {
	case string:
		{
			_log.Panic(v, fields...)
		}
	case error:
		{
			_log.Panic(v.Error(), fields...)
		}
	}
}

func Fatal(msg interface{}, fields ...zap.Field) {
	switch v := msg.(type) {
	case string:
		{
			_log.Fatal(v, fields...)
		}
	case error:
		{
			_log.Fatal(v.Error(), fields...)
		}
	}
}
