package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

// 重写gormLogger逻辑
var (
	infoStr     = "\r\n事件:%s\r\n%s\r\n"
	warnStr     = "\r\n事件:SQL执行,[warn]%s\r\n"
	errStr      = "\r\n事件:SQL执行,[err]%s\r\n"
	traceStr    = "\r\n事件:SQL执行\r\n时间:%.3fms\r\n行数:%v\r\nSQL:%s\r\n%s\r\n"
	traceErrStr = "\r\n事件:SQL执行\r\nErr:%s\r\n时间:%.3fms\r\n行数:%v\r\nSQL：%s\r\n%s\r\n"
)

type GormCustomLogger struct {
	logLevel logger.LogLevel
	mLog     *zap.SugaredLogger
}

func NewGormCustomLogger() *GormCustomLogger {
	log := zap.L().Sugar()
	return &GormCustomLogger{
		mLog:     log,
		logLevel: logger.Info,
	}
}
func (g *GormCustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLog := *g
	newLog.logLevel = level
	return &newLog
}

func (g *GormCustomLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if g.logLevel >= logger.Info {
		str := fmt.Sprintf(s, i...)
		g.mLog.Infof(infoStr, str, LoggerEndStr)
	}
}

func (g *GormCustomLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if g.logLevel >= logger.Warn {
		g.mLog.Warnf(warnStr+s, append([]interface{}{LoggerEndStr}, i...)...)
	}
}

func (g *GormCustomLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if g.logLevel >= logger.Error {
		g.mLog.Errorf(errStr+s, i...)
	}
}

func (g *GormCustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.logLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && g.logLevel >= logger.Error:
		g.mLog.Infof(traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql, LoggerEndStr)
	case g.logLevel >= logger.Info:
		g.mLog.Infof(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql, LoggerEndStr)
	}
}
