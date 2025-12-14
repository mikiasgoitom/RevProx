package contract

import (
	"context"

	valueobject "github.com/mikiasgoitom/RevProx/internal/domain/valueObject"
)

type ILogger interface {
	Info(ctx context.Context, msg string, fields ...valueobject.LogField)
	Debug(ctx context.Context, msg string, fields ...valueobject.LogField)
	Warn(ctx context.Context, msg string, fields ...valueobject.LogField)
	Error(ctx context.Context, msg string, fields ...valueobject.LogField)
	Fatal(ctx context.Context, msg string, fields ...valueobject.LogField)
}