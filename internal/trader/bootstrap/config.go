package bootstrap

import "os"

type Config struct {
	Database
	Sentry
	StatisticoDataService
	StatisticoOddsWarehouseService
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

func BuildConfig() *Config {
	config := Config{}

	config.Database = Database{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	config.Sentry = Sentry{DSN: os.Getenv("SENTRY_DSN")}

	config.StatisticoDataService = StatisticoDataService{
		Host: os.Getenv("STATISTICO_DATA_SERVICE_HOST"),
		Port: os.Getenv("STATISTICO_DATA_SERVICE_PORT"),
	}

	config.StatisticoOddsWarehouseService = StatisticoOddsWarehouseService{
		Host: os.Getenv("STATISTICO_ODDS_WAREHOUSE_SERVICE_HOST"),
		Port: os.Getenv("STATISTICO_ODDS_WAREHOUSE_SERVICE_PORT"),
	}

	return &config
}
