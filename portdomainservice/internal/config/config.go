package config

import (
	"log"
	"strings"

	"github.com/pkg/errors"
	v "github.com/spf13/viper"
)

type Config struct {
	GRPCHost     string `mapstructure:"grpc_host"`
	GRPCPort     string `mapstructure:"grpc_port"`
	DBConnString string `mapstructure:"db_conn_url"`
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
