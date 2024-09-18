package main

import (
    "io/fs"
    "log/slog"

    "github.com/spf13/viper"
)

type Config struct {
    Address string
    Secret  string
    PGUrl   string
}

func (cfg Config) LogValue() slog.Value {
    return slog.GroupValue(
        slog.String("address", cfg.Address),
        slog.String("secret", "*"),
        slog.String("DATABASE_URL", cfg.PGUrl),
    )
}

func setDefaults() {
    viper.SetDefault("HOST", "0.0.0.0")
    viper.SetDefault("PORT", "8080")
    viper.SetDefault("SECRET", "secret")
    //viper.SetDefault("DATABASE_URL", "host=0.0.0.0 user=postgres password=postgres dbname=postgres sslmode=disable")
    viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable")
}

func makeConfig() Config {
    return Config{
        Address: viper.GetString("HOST") + ":" + viper.GetString("PORT"),
        Secret:  viper.GetString("SECRET"),
        PGUrl:   viper.GetString("DATABASE_URL"),
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
