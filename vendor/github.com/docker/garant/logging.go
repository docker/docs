package garant

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"github.com/docker/distribution/context"
	"github.com/docker/garant/config"
)

// configureLogging prepares a new context with a logger using the
// configuration.
func configureLoggingContext(loggingOpts config.Logging) context.Context {
	level, err := log.ParseLevel(loggingOpts.Level)
	if err != nil {
		level = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q	", loggingOpts.Level, err, level)
	}

	var formatter log.Formatter

	switch loggingOpts.Formatter {
	case "", "text":
		formatter = &log.TextFormatter{}
	case "json":
		formatter = &log.JSONFormatter{}
	case "logstash":
		formatter = &logstash.LogstashFormatter{}
	default:
		formatter = &log.TextFormatter{}
		log.Warnf("unsupported logging formatter %q, using default text formatter", loggingOpts.Formatter)
	}

	log.SetLevel(level)
	log.SetFormatter(formatter)

	ctx := context.Background()

	if len(loggingOpts.Fields) > 0 {
		// build up the static fields, if present.
		fields := make(map[string]interface{})
		var keys []interface{}
		for k, v := range loggingOpts.Fields {
			fields[fmt.Sprint(k)] = v
			keys = append(keys, k)
		}
		ctx = context.WithValues(ctx, fields)
		ctx = context.WithLogger(ctx, context.GetLogger(ctx, keys...))
	}

	return ctx
}
