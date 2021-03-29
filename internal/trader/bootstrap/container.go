package bootstrap

import (
	"database/sql"
	"fmt"
	"github.com/evalphobia/logrus_sentry"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Container struct {
	Clock  clockwork.Clock
	Config *Config
	Database *sql.DB
	Logger *logrus.Logger
}

func BuildContainer(config *Config) Container {
	c := Container{
		Config: config,
	}

	c.Clock = clockwork.NewRealClock()
	c.Database = databaseConnection(config)
	c.Logger = logger(config)

	return c
}

func databaseConnection(config *Config) *sql.DB {
	db := config.Database

	dsn := "host=%s port=%s user=%s " +
		"password=%s dbname=%s sslmode=disable"

	psqlInfo := fmt.Sprintf(dsn, db.Host, db.Port, db.User, db.Password, db.Name)

	conn, err := sql.Open(db.Driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(25)

	return conn
}

func logger(config *Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	tags := map[string]string{
		"application": "statistico-strategy",
	}

	levels := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}

	hook, err := logrus_sentry.NewWithTagsSentryHook(config.Sentry.DSN, tags, levels)

	if err == nil {
		hook.Timeout = 20 * time.Second
		hook.StacktraceConfiguration.Enable = true
		hook.StacktraceConfiguration.IncludeErrorBreadcrumb = true
		logger.AddHook(hook)
	}

	return logger
}
