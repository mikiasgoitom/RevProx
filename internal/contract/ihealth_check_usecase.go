package contract
import "context"

type IHealthCheckUseCase interface {
	Readiness(ctx context.Context) error
	Liveness(ctx context.Context) error
}