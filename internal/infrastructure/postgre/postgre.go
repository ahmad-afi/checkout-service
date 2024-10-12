package postgre

import (
	"checkout-service/internal/utils"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host      string `env:"POSTGRES_HOST"`
	Port      int    `env:"POSTGRES_PORT"`
	Ordername string `env:"POSTGRES_USER"`
	Password  string `env:"POSTGRES_PASSWORD"`
	DbName    string `env:"POSTGRES_DB"`

	// 	MaxLifetime        int    `env:"MYSQL_maxLifetime"`
	// 	MinIdleConnections int    `env:"MYSQL_minIdleConnections"`
	// 	MaxOpenConnections int    `env:"MYSQL_maxOpenConnections"`
}

func Init() (*sqlx.DB, error) {
	var db *sqlx.DB

	var cfg = Config{
		Host:      utils.EnvString("POSTGRES_HOST"),
		Port:      utils.EnvInt("POSTGRES_PORT"),
		Ordername: utils.EnvString("POSTGRES_USER"),
		Password:  utils.EnvString("POSTGRES_PASSWORD"),
		DbName:    utils.EnvString("POSTGRES_DB"),
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Ordername, cfg.Password, cfg.DbName)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	err = db.Ping()

	if err != nil {
		return db, err
	}

	return db, nil
}
