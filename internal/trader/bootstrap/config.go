package bootstrap

import (
	"net/http"
	"os"
)

type Config struct {
	AWS
	Database
	HTTPClient  *http.Client
	QueueDriver string
	Sentry
	StatisticoDataService
	StatisticoOddsWarehouseService
	User
}

type AWS struct {
	Key      string
	Region  string
	CognitoUserPoolID string
	QueueUrl string
	Secret   string
}

type Database struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Sentry struct {
	DSN string
}

type StatisticoDataService struct {
	Host string
	Port string
}

type StatisticoOddsWarehouseService struct {
	Host string
	Port string
}

// User is a temporary struct containing a hardcoded user until abstract user management functionality is implemented
type User struct {
	ID    string
	Email string
	BetFairUserName string
	BetFairPassword string
	BetFairKey string
}

func BuildConfig() *Config {
	config := Config{}

	config.AWS = AWS{
		CognitoUserPoolID: os.Getenv("AWS_USER_POOL_ID"),
		Key:               os.Getenv("AWS_KEY"),
		QueueUrl:          os.Getenv("AWS_QUEUE_URL"),
		Region:            os.Getenv("AWS_REGION"),
		Secret:            os.Getenv("AWS_SECRET"),
	}

	config.Database = Database{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	config.HTTPClient = &http.Client{}

	config.QueueDriver = os.Getenv("QUEUE_DRIVER")

	config.Sentry = Sentry{DSN: os.Getenv("SENTRY_DSN")}

	config.StatisticoDataService = StatisticoDataService{
		Host: os.Getenv("STATISTICO_DATA_SERVICE_HOST"),
		Port: os.Getenv("STATISTICO_DATA_SERVICE_PORT"),
	}

	config.StatisticoOddsWarehouseService = StatisticoOddsWarehouseService{
		Host: os.Getenv("STATISTICO_ODDS_WAREHOUSE_SERVICE_HOST"),
		Port: os.Getenv("STATISTICO_ODDS_WAREHOUSE_SERVICE_PORT"),
	}

	config.User = User{
		ID:              os.Getenv("USER_ID"),
		Email:           os.Getenv("USER_EMAIL_ADDRESS"),
		BetFairUserName: os.Getenv("BETFAIR_USERNAME"),
		BetFairPassword: os.Getenv("BETFAIR_PASSWORD"),
		BetFairKey:      os.Getenv("BETFAIR_KEY"),
	}

	return &config
}
