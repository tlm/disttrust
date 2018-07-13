package log

// Thanks to https://github.com/containerd/containerd/blob/master/log/context.go
// for the idea

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})

	if logger == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return logger.(*logrus.Entry)
}

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}
