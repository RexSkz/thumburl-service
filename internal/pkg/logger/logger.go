package logger

import (
	"context"
	"thumburl-service/internal/middleware"

	"go.uber.org/zap"
)

var _l, _ = zap.NewProduction()
var l = _l.WithOptions(zap.AddCallerSkip(1))

func Debugw(ctx context.Context, s string, v ...interface{}) {
	defer l.Sync()
	l.Sugar().Debugw(s, withContext(ctx, v)...)
}

func Infow(ctx context.Context, s string, v ...interface{}) {
	defer l.Sync()
	l.Sugar().Infow(s, withContext(ctx, v)...)
}

func Warnw(ctx context.Context, s string, v ...interface{}) {
	defer l.Sync()
	l.Sugar().Warnw(s, withContext(ctx, v)...)
}

func Errorw(ctx context.Context, s string, v ...interface{}) {
	defer l.Sync()
	l.Sugar().Errorw(s, withContext(ctx, v)...)
}

func Panicw(ctx context.Context, s string, v ...interface{}) {
	defer l.Sync()
	l.Sugar().Panicw(s, withContext(ctx, v)...)
}

func withContext(ctx context.Context, v []interface{}) []interface{} {
	vv := make([]interface{}, len(v)+2)
	vv[0] = middleware.TraceIDKey
	vv[1] = ctx.Value(middleware.TraceIDKey)
	if vv[1] == nil {
		vv[1] = middleware.TraceIDGlobal
	}
	for i, val := range v {
		vv[i+2] = val
	}
	return vv
}
