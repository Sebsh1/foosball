package main

import (
	"context"
	"fmt"
	"foosball/internal/authentication"
	"foosball/internal/connectors/mysql"
	"foosball/internal/invite"
	"foosball/internal/match"
	"foosball/internal/organization"
	"foosball/internal/rating"
	"foosball/internal/rest"
	"foosball/internal/user"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	ctx := context.Background()

	config, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	log := getLogger(config.Log)

	db, err := mysql.NewClient(ctx, config.DB)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}

	if err := db.AutoMigrate(
		&match.Match{},
		&match.Set{},
		&organization.Organization{},
		&user.User{},
		&invite.Invite{},
	); err != nil {
		log.WithError(err).Fatal("failed to auto migrate database")
	}

	userReposotory := user.NewRepository(db)
	userService := user.NewService(userReposotory)

	ratingService := rating.NewService(config.Rating, userService)

	organizationRepository := organization.NewRepository(db)
	organizationService := organization.NewService(organizationRepository)

	authService := authentication.NewService(config.Auth, userService)

	inviteRepo := invite.NewRepository(db)
	inviteService := invite.NewService(inviteRepo, userService, organizationService)

	matchRepo := match.NewRepository(db)
	matchService := match.NewService(matchRepo)

	httpServer, err := rest.NewServer(
		config.Rest,
		log,
		authService,
		userService,
		organizationService,
		inviteService,
		matchService,
		ratingService,
	)
	if err != nil {
		log.WithError(err).Fatal("failed to create http server")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.WithField("port", config.Rest.Port).Info("REST server starting")
	go func() {
		if err := httpServer.Start(); err != nil {
			log.WithError(err).Error("failed to start http server")
			cancel()
		}
	}()

	log.Info("Ready")

	<-ctx.Done()
	log.Info("shutting down")
	shutdownCtx, stop := context.WithTimeout(context.Background(), shutdownPeriod)
	defer stop()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("failed to shutdown internal http server")
	}

}

var cfgFile string

const shutdownPeriod = 15 * time.Second

type LogConfig struct {
	Level string `mapstructure:"level" validate:"required"`
}

type Config struct {
	Log    LogConfig             `mapstructure:"log" validate:"dive"`
	Rest   rest.Config           `mapstructure:"rest" validate:"dive"`
	Auth   authentication.Config `mapstructure:"auth" validate:"dive"`
	DB     mysql.Config          `mapstructure:"db" validate:"dive"`
	Rating rating.Config         `mapstructure:"rating" validate:"dive"`
}

func initConfig() {
	viper.AddConfigPath(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()
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

func getLogger(config LogConfig) *logrus.Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(config.Level)
	if err != nil {
		lvl = logrus.InfoLevel
		logger.WithError(err).Warnf("failed to parse log level, setting log level to '%s'", lvl)
	}
	logger.SetLevel(lvl)
	return logger
}

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
		switch fieldv.Kind() {
		case reflect.Struct:
			bindEnvs(fieldv.Interface(), parts...)
		default:
			viper.BindEnv(strings.Join(parts, "."))
		}
	}
}
