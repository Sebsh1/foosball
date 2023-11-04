package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const shutdownPeriod = 15 * time.Second

type Config struct {
	LogEnv        string        `mapstructure:"LOG_ENV"`
	DBDSN         string        `mapstructure:"DB_DSN"`
	Port          int           `mapstructure:"PORT"`
	JWTSecret     string        `mapstructure:"JWT_SECRET"`
	JWTExpiration time.Duration `mapstructure:"JWT_EXPIRATION"`
}

var rootCmd = &cobra.Command{
	Use:   "matchlog",
	Short: "Matchlog is a service for tracking match results",
	Long: "Matchlog is a service for tracking match results. " +
		"It provides a REST API for creating and managing users, Clubs, matches, and more.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("Failed to execute root command", zap.Error(err))
	}
}

func init() { //nolint:gochecknoinits
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Fatal("Can't find the file .env", zap.Error(err))
	}
}

func loadConfig() *Config {
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		zap.L().Fatal("Environment can't be loaded", zap.Error(err))
	}

	return &config
}

func GetLogger(logEnv string) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	switch logEnv {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewProduction()
	}
	if err != nil {
		zap.L().Fatal("Failed to build logger", zap.Error(err))
	}

	logger.Info("Logger initialized",
		zap.String("logEnv", logEnv),
	)

	return logger.Sugar()
}
