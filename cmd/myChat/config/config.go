package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"runtime"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	FrontUrl   string `mapstructure:"FRONT_URL"`
	DebugMode  bool   `mapstructure:"DEBUG"`
	DbName     string `mapstructure:"DB_NAME"`
	DbUrl      string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (config *Config, err error) {

	pathConfig := getAbsolutePath()
	env := getenv("ENV")

	configName := fmt.Sprintf("%s_config", env)

	viper.AddConfigPath(pathConfig)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)

	return
}

func getAbsolutePath() string {
	_, filename, _, _ := runtime.Caller(0)

	return path.Join(path.Dir(filename), "../../../")
}

func getenv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "local"
	}
	return value
}
