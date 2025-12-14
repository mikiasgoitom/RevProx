package logger

import (
	"context"

	"github.com/mikiasgoitom/RevProx/internal/contract"
	valueobject "github.com/mikiasgoitom/RevProx/internal/domain/valueObject"
	"go.uber.org/zap"
)

type ZapAdapter struct {
	logger *zap.Logger
}

func NewZapAdapter(isProduction bool) (contract.ILogger, error) {
	var zapLogger *zap.Logger
	var err error

	if isProduction {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}
	return &ZapAdapter{logger: zapLogger}, nil
}

func (z *ZapAdapter) toZapFields(fields ...valueobject.LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}	
	return zapFields
}

func (z *ZapAdapter) Info(ctx context.Context, msg string, fields ...valueobject.LogField){
	z.logger.Info(msg, z.toZapFields(fields...)...)
}
func (z *ZapAdapter) Debug(ctx context.Context, msg string, fields ...valueobject.LogField){
	z.logger.Debug(msg, z.toZapFields(fields...)...)
}

func (z *ZapAdapter) Warn(ctx context.Context, msg string, fields ...valueobject.LogField){
	z.logger.Warn(msg, z.toZapFields(fields...)...)
}
func (z *ZapAdapter) Error(ctx context.Context, msg string, fields ...valueobject.LogField){
	z.logger.Error(msg, z.toZapFields(fields...)...)
}
func (z *ZapAdapter) Fatal(ctx context.Context, msg string, fields ...valueobject.LogField){
	z.logger.Fatal(msg, z.toZapFields(fields...)...)
}