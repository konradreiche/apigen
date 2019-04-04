package api

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggingMiddleware struct {
	log *logrus.Logger
	api API
}

func NewLoggingMiddleware(api API, log *logrus.Logger) API {
	return &loggingMiddleware{
		log: log,
		api: api,
	}
}

func (l *loggingMiddleware) GetPrice(ctx context.Context, req GetPriceRequest) (*GetPriceResponse, error) {
	var err error
	defer func() {
		fields := logrus.Fields{
			"method":     "GetPrice",
			"user_agent": ctx.Value("userAgent"),
		}
		if err == nil {
			l.log.WithFields(fields).Info("request")
		} else {
			fields["error"] = err
			l.log.WithFields(fields).Error("request failed")
		}
	}()
	resp, err := l.api.GetPrice(ctx, req)
	return resp, err
}
