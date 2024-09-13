package main

import (
	"io/fs"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Address string
	Secret  string
}

func (cfg Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("address", cfg.Address),
		slog.String("secret", "*"),
	)
}

func setDefaults() {
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("SECRET", "secret")
}

func makeConfig() Config {
	return Config{
		Address: viper.GetString("HOST") + ":" + viper.GetString("PORT"),
		Secret:  viper.GetString("SECRET"),
	}
}

func loadConfig() (Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetEnvPrefix("API")
	setDefaults()
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError: // Do nothing
		case *fs.PathError: // Do nothing
		default:
			return Config{}, err
		}
	}
	return makeConfig(), nil
}
