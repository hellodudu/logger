package logger

import (
	"errors"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// 初始化logger文件名字，是否采集日志到loki，如果是跑在docker中的话不需要开启loki
	Init("game server", true, "")
	EnableLoki(true)
	//SetLokiConfig("http://localhost:3100/api/prom/push", 1024, 5)
	Info("test")
	Warn("warn")
	Error("error")

	err := errors.New("error 404 found")
	fields := map[string]interface{}{
		"error": err,
		"url":   "http://google.com",
	}
	WithFieldsWarn(fields, "ping to google")

	tm := time.Now()
	md := map[string]string{"log_level": "info", "access_level": "root"}
	fields2 := map[string]interface{}{
		"time": tm,
	}
	WithFieldsInfo(fields2, md)

	time.Sleep(time.Second * 5)
}
