package configservice

import (
	"strings"

	"github.com/mikiasgoitom/RevProx/internal/config"
	"github.com/mikiasgoitom/RevProx/internal/contract"
	"github.com/spf13/viper"
)

type ViperAdapter struct {
}

func NewViperAdapter() contract.IConfigService {
	return &ViperAdapter{}
}

func (v *ViperAdapter) Load() (config.Config, error) {
	var cfg config.Config
	// system default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("cache.max_cost", "100MB")
	viper.SetDefault("cache.num_counters", 1e6)
	viper.SetDefault("origin.base_url", "http://localhost:3000")

	// read from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	
	// read environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return cfg, err
		}
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}