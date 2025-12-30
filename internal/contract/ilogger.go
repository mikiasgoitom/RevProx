package contract

import (
	"context"
	"github.com/mikiasgoitom/RevProx/internal/domain/valueobject"
)

type ILogger interface {
	Info(ctx context.Context, msg string, fields ...valueobject.LogField)
	Debug(ctx context.Context, msg string, fields ...valueobject.LogField)
	Warn(ctx context.Context, msg string, fields ...valueobject.LogField)
	Error(ctx context.Context, msg string, fields ...valueobject.LogField)
	Fatal(ctx context.Context, msg string, fields ...valueobject.LogField)
}