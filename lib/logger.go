package lib

import "go.uber.org/zap"

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
}

type noOpLogger struct{}

func (noOpLogger) Debugf(template string, args ...interface{}) {}
func (noOpLogger) Infof(template string, args ...interface{})  {}
func (noOpLogger) Warnf(template string, args ...interface{})  {}
func (noOpLogger) Errorf(template string, args ...interface{}) {}
func (noOpLogger) Fatalf(template string, args ...interface{}) {}
func (noOpLogger) Panicf(template string, args ...interface{}) {}

type ZapAdapter struct {
	Logger *zap.SugaredLogger
}

func (a ZapAdapter) Debugf(template string, args ...interface{}) {
	a.Logger.Debugf(template, args...)
}

func (a ZapAdapter) Infof(template string, args ...interface{}) {
	a.Logger.Infof(template, args...)
}

func (a ZapAdapter) Warnf(template string, args ...interface{}) {
	a.Logger.Warnf(template, args...)
}

func (a ZapAdapter) Errorf(template string, args ...interface{}) {
	a.Logger.Errorf(template, args...)
}

func (a ZapAdapter) Fatalf(template string, args ...interface{}) {
	a.Logger.Fatalf(template, args...)
}

func (a ZapAdapter) Panicf(template string, args ...interface{}) {
	a.Logger.Panicf(template, args...)
}
