package kinesis

import (
	"go.uber.org/zap"
)

type LoggerAdapter struct {
	logger *zap.SugaredLogger
}

func (l *LoggerAdapter) Log(args ...interface{}) {
	l.logger.Debug(args...)
}
