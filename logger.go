package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	logrusloki "github.com/hellodudu/logrus-loki"
	logrus "github.com/sirupsen/logrus"
)

var (
	lokiEndPoint  = "http://192.168.1.244:3100/api/prom/push"
	lokiBatchSize = 1024
	lokiBatchWait = 3

	std *Logger = nil
)

type Logger struct {
	*logrus.Logger
	sync.Mutex

	loki       *logrusloki.Loki
	enableLoki bool
}

// please disable loki when run in docker container, it is better to choose promtail gathering container's log
func Init(fn string, enableLoki bool, lokiUrl string) {
	if std != nil && std.loki != nil {
		std.loki.Close()
	}

	std = newLogger(fn, enableLoki, lokiUrl)
}

// default settings
func newLogger(fn string, enableLoki bool, lokiUrl string) *Logger {
	l := &Logger{
		Logger:     logrus.StandardLogger(),
		enableLoki: enableLoki,
	}

	// log file name
	t := time.Now()
	fileTime := fmt.Sprintf("%d-%d-%d %d-%d-%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	logFn := fmt.Sprintf("log/%s_%s.log", fileTime, fn)

	file, err := os.OpenFile(logFn, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	// set writer
	l.SetOutput(io.MultiWriter(os.Stdout, file))

	// report caller
	l.SetReportCaller(true)

	// set formatter
	fileFormatter := new(logrus.TextFormatter)
	fileFormatter.TimestampFormat = "2006-01-02 15:04:05"
	fileFormatter.FullTimestamp = true
	l.SetFormatter(fileFormatter)

	// loki
	if enableLoki {
		if len(lokiUrl) == 0 {
			lokiUrl = lokiEndPoint
		}

		l.loki, err = logrusloki.NewLoki(lokiUrl, lokiBatchSize, lokiBatchWait)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to new loki, %v\n", err)
		}
	}

	return l
}

func (l *Logger) log(level logrus.Level, args ...interface{}) {
	entry := logrus.NewEntry(std.Logger)

	if l.enableLoki && l.loki != nil {
		entry.Message = fmt.Sprint(args...)
		entry.Level = level
		entry.Time = time.Now()
		l.loki.Fire(entry)
	}

	entry.Log(level, args...)
}

func (l *Logger) logFields(level logrus.Level, fields map[string]interface{}, args ...interface{}) {
	entry := l.WithFields(fields)

	if l.enableLoki && l.loki != nil {
		entry.Message = fmt.Sprint(args...)
		entry.Data = fields
		entry.Level = level
		entry.Time = time.Now()
		l.loki.Fire(entry)
	}

	entry.Log(level, args...)
}

func EnableLoki(enable bool) {
	std.Lock()
	defer std.Unlock()

	std.enableLoki = enable
}

func SetLokiConfig(url string, batchSize int, batchWait int) error {
	std.Lock()
	defer std.Unlock()

	if std.loki != nil {
		std.loki.Close()
	}

	var err error
	std.loki, err = logrusloki.NewLoki(url, batchSize, batchWait)
	if err != nil {
		return fmt.Errorf("SetLokiConfig failed: %w", err)
	}

	return nil
}

func Trace(args ...interface{}) {
	std.log(logrus.TraceLevel, args...)
}

func Debug(args ...interface{}) {
	std.log(logrus.DebugLevel, args...)
}

func Info(args ...interface{}) {
	std.log(logrus.InfoLevel, args...)
}

func Warn(args ...interface{}) {
	std.log(logrus.WarnLevel, args...)
}

func Error(args ...interface{}) {
	std.log(logrus.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	std.log(logrus.FatalLevel, args...)
}

func Panic(args ...interface{}) {
	std.log(logrus.PanicLevel, args...)
}

func WithFieldsTrace(fields map[string]interface{}, args ...interface{}) {
	std.logFields(logrus.TraceLevel, fields, args...)
}

func WithFieldsDebug(fields map[string]interface{}, args ...interface{}) {
	std.logFields(logrus.DebugLevel, fields, args...)
}

func WithFieldsInfo(fields map[string]interface{}, args ...interface{}) {
	std.logFields(logrus.InfoLevel, fields, args...)
}

func WithFieldsWarn(fields map[string]interface{}, args ...interface{}) {
	std.logFields(logrus.WarnLevel, fields, args...)
}

func WithFieldsError(fields map[string]interface{}, args ...interface{}) {
	std.logFields(logrus.ErrorLevel, fields, args...)
}

func WithFieldsFatal(fields map[string]interface{}, args ...interface{}) {
	std.logFields(logrus.FatalLevel, fields, args...)
}

func WithFieldsPanic(fields map[string]interface{}, args ...interface{}) {
	std.logFields(logrus.PanicLevel, fields, args...)
}
