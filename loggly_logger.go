package gomol_loggly

import (
	"fmt"
	"time"

	"github.com/aphistic/gomol"
	"github.com/segmentio/go-loggly"
)

var logglyClients map[string]*loggly.Client

func init() {
	logglyClients = make(map[string]*loggly.Client, 0)
}

type LogglyLoggerConfig struct {
	Token string
}

func NewLogglyLoggerConfig() *LogglyLoggerConfig {
	return &LogglyLoggerConfig{}
}

type LogglyLogger struct {
	base          *gomol.Base
	isInitialized bool
	config        *LogglyLoggerConfig
}

func NewLogglyLogger(cfg *LogglyLoggerConfig) (*LogglyLogger, error) {
	l := &LogglyLogger{
		config: cfg,
	}
	return l, nil
}

func (l *LogglyLogger) getClient() *loggly.Client {
	c := logglyClients[l.config.Token]
	if c == nil {
		c = loggly.New(l.config.Token)
		c.Level = loggly.DEBUG
		logglyClients[l.config.Token] = c
	}
	return c
}

func (l *LogglyLogger) getFacility(m map[string]interface{}) string {
	if m != nil {
		if facility, ok := m["facility"]; ok {
			return fmt.Sprint(facility)
		}
	}
	if facility := l.base.BaseAttrs.GetAttr("facility"); facility != nil {
		return fmt.Sprint(facility)
	}
	return ""
}

func (l *LogglyLogger) getMsg(m map[string]interface{}, msg string, a ...interface{}) loggly.Message {
	lm := loggly.Message{"message": fmt.Sprintf(msg, a...)}

	for key, val := range l.base.BaseAttrs.Attrs() {
		lm[key] = val
	}
	if m != nil {
		for key, val := range m {
			lm[key] = fmt.Sprintf("%v", val)
		}
	}

	return lm
}

func (l *LogglyLogger) SetBase(base *gomol.Base) {
	l.base = base
}

func (l *LogglyLogger) IsInitialized() bool {
	return l.isInitialized
}

func (l *LogglyLogger) InitLogger() error {
	l.isInitialized = true
	return nil
}

func (l *LogglyLogger) ShutdownLogger() error {
	c := logglyClients[l.config.Token]
	if c != nil {
		err := c.Flush()
		if err != nil {
			return err
		}

		delete(logglyClients, l.config.Token)
	}

	l.isInitialized = false

	return nil
}

func (l *LogglyLogger) Logm(timestamp time.Time, level gomol.LogLevel, m map[string]interface{}, msg string) error {
	lm := l.getMsg(nil, msg)
	switch level {
	case gomol.LevelDebug:
		l.getClient().Debug(l.getFacility(nil), lm)
	case gomol.LevelInfo:
		l.getClient().Info(l.getFacility(nil), lm)
	case gomol.LevelWarning:
		l.getClient().Warn(l.getFacility(nil), lm)
	case gomol.LevelError:
		l.getClient().Error(l.getFacility(nil), lm)
	case gomol.LevelFatal:
		l.getClient().Critical(l.getFacility(nil), lm)
	}
	return nil
}
