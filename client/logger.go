package client

type Logger interface {
	Error(c *QQClient, msg string, args ...interface{})
	Warning(c *QQClient, msg string, args ...interface{})
	Info(c *QQClient, msg string, args ...interface{})
	Debug(c *QQClient, msg string, args ...interface{})
	Trace(c *QQClient, msg string, args ...interface{})
}

type nopLogger struct{}

func (l nopLogger) Error(c *QQClient, msg string, args ...interface{})   {}
func (l nopLogger) Warning(c *QQClient, msg string, args ...interface{}) {}
func (l nopLogger) Info(c *QQClient, msg string, args ...interface{})    {}
func (l nopLogger) Debug(c *QQClient, msg string, args ...interface{})   {}
func (l nopLogger) Trace(c *QQClient, msg string, args ...interface{})   {}
