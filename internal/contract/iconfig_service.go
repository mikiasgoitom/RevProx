package contract

import "github.com/mikiasgoitom/RevProx/internal/config"

type IConfigService interface {
	Load() (config.Config, error)
}