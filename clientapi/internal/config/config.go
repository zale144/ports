package config

import (
	"log"
	"strings"

	"github.com/pkg/errors"
	v "github.com/spf13/viper"
)

type Config struct {
	HTTPPort  int64  `mapstructure:"http_port"`
	BatchSize int64  `mapstructure:"batch_size"`
	GRPCURL   string `mapstructure:"grpc_url"`
}

func Configure() (*Config, error) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AddConfigPath("./")
	v.SetConfigName("config")
	if err := v.MergeInConfig(); err != nil {
		log.Println(err)
	}
	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "could not configure")
	}

	return &config, nil
}
