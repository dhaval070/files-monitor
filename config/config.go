package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DataDir string `mapstructure:"data_dir" validate:"required"`
	DB_DSN  string `mapstructure:"db_dsn" validate:"required"`
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
}

func MustReadConfig() *Config {
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("error while reading config file ", err)
	}

	var c = Config{}
	err = viper.Unmarshal(&c)

	if err != nil {
		log.Fatal("error while unmarshalling config file ", err)
	}

	v := validator.New()
	if err = v.Struct(&c); err != nil {
		log.Fatal(err)
	}

	return &c
}
