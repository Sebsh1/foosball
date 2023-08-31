package cmd

import (
	"encoding/json"
	"fmt"
	"matchlog/internal/authentication"
	"matchlog/internal/rest"
	"matchlog/pkg/database"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string

const shutdownPeriod = 15 * time.Second

type Config struct {
	Log            LogConfig             `mapstructure:"log" validate:"dive"`
	Rest           rest.Config           `mapstructure:"rest" validate:"dive"`
	DB             database.Config       `mapstructure:"db" validate:"dive"`
	Authentication authentication.Config `mapstructure:"authentication" validate:"dive"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

var rootCmd = &cobra.Command{
	Use:   "matchlog",
	Short: "Matchlog is a service for tracking match results",
	Long: "Matchlog is a service for tracking match results. " +
		"It provides a REST API for creating and managing users, organizations, matches, and more.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("Failed to execute root command", zap.Error(err))
	}
}

func init() { //nolint:gochecknoinits
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err == nil {
		zap.L().Info("Using config", zap.String("file", viper.ConfigFileUsed()))
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func loadConfig(configs ...string) (*Config, error) {
	var config Config
	defaults.SetDefaults(&config)
	bindEnvs(config)

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	match := regexp.MustCompile(`.*`)
	if len(configs) != 0 {
		match = regexp.MustCompile(strings.ToLower(fmt.Sprintf("^Config.(%s)", strings.Join(configs, "|"))))
	}

	err = validator.New().StructFiltered(config, func(ns []byte) bool {
		return !match.MatchString(strings.ToLower(string(ns)))
	})
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetLogger(config LogConfig) *zap.SugaredLogger {
	configString := fmt.Sprintf(`{
		"level": "%s",
		"encoding": "%s",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`, config.Level, config.Encoding)
	var cfg zap.Config
	if err := json.Unmarshal([]byte(configString), &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

// Adapted from https://github.com/spf13/viper/issues/188#issuecomment-401431526
func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)
		name := strings.ToLower(t.Name)
		tag, ok := t.Tag.Lookup("mapstructure")
		if ok {
			name = tag
		}
		parts := append(parts, name)
		switch fieldv.Kind() { //nolint:exhaustive
		case reflect.Struct:
			bindEnvs(fieldv.Interface(), parts...)
		default:
			if err := viper.BindEnv(strings.Join(parts, ".")); err != nil {
				zap.L().Fatal("Failed to bind environment variable", zap.String("variable", strings.Join(parts, ".")), zap.Error(err))
			}
		}
	}
}
