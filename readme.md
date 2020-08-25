Logger combines a wrapper of logrus and schoentoon/logrus-loki, not only use stdout as standard output, but alse add a new file writer output. It can easily sending entries of logrus to loki without docker container.

## Usage

```golang
    Init("game server", true, "")
	EnableLoki(true)
	SetLokiConfig("http://localhost:3100/api/prom/push", 1024, 5)
	Info("test")
	Warn("warn")
	Error("error")

	err := errors.New("error 404 found")
	fields := map[string]interface{}{
		"error": err,
		"url":   "http://google.com",
	}
	WithFieldsWarn(fields, "ping to google")

```

