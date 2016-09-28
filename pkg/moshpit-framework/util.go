package moshpit

import (
	"github.com/Sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/net/context"
)

func LoggerFromCtx(ctx context.Context) logrus.FieldLogger {
	return ctx.Value("logger").(logrus.FieldLogger)
}

func MakeCtx(member string, level logrus.Level) context.Context {
	// root context and logger
	var logger logrus.FieldLogger
	loggr := logrus.New()
	loggr.Formatter = new(prefixed.TextFormatter)
	loggr.Level = level
	logger = loggr
	logger = logger.WithField("prefix", member)
	return context.WithValue(context.Background(), "logger", logger)
}
