package integration_test

import (
	"checkout-service/internal/helper"
	"checkout-service/internal/infrastructure/postgre"
	"checkout-service/internal/utils"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sqlx.DB
var pool *dockertest.Pool
var resource *dockertest.Resource

func InitDocketTest() {

	err := godotenv.Load(fmt.Sprintf("%s/%s", helper.ProjectRootPath, ".env.test"))
	if err != nil {
		panic(err)
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	var cfg = postgre.Config{
		Host:     utils.EnvString("POSTGRES_HOST"),
		Port:     utils.EnvInt("POSTGRES_PORT"),
		Username: utils.EnvString("POSTGRES_USER"),
		Password: utils.EnvString("POSTGRES_PASSWORD"),
		DbName:   utils.EnvString("POSTGRES_DB"),
	}

	var envString = []string{
		"POSTGRES_PASSWORD=" + cfg.Password,
		"POSTGRES_USER=" + cfg.Username,
		"POSTGRES_DB=" + cfg.DbName,
		"listen_addresses = '*'",

		// "POSTGRES_PASSWORD=secret",
		// "POSTGRES_USER=user_name",
		// "POSTGRES_DB=dbname",
		// "listen_addresses = '*'",
	}

	// pulls an image, creates a container based on it and runs it
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "15.4",
		Env:          envString,
		ExposedPorts: []string{fmt.Sprint(cfg.Port)},
		Name:         "pg-test",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "127.0.0.1", HostPort: fmt.Sprint(cfg.Port) + "/tcp"}},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// hostAndPort := resource.GetHostPort(fmt.Sprint(cfg.Port) + "/tcp")
	// databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.DbName)

	// hostAndPort := resource.GetHostPort("5432/tcp")
	// databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName)

	log.Println("Connecting to database on dsn: ", dsn)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	if err = pool.Retry(func() error {
		db, err = sqlx.Open("postgres", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// defer func() {
	// 	if err := pool.Purge(resource); err != nil {
	// 		log.Fatalf("Could not purge resource: %s", err)
	// 	}
	// }()

	log.Println("postgres container created")
	// run tests
	RunMigration(db)
}

func ClearingDocker() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func RunMigration(db *sqlx.DB) (err error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		DatabaseName: utils.EnvString("POSTGRES_DB"),
	})
	if err != nil {
		helper.Logger(helper.LoggerLevelError, "error connect db", err)
		return
	}
	defer func() {
		fmt.Println("closing connection")
		driver.Close()
	}()

	m, err := migrate.NewWithDatabaseInstance("file://../../migrations", "postgres", driver)
	if err != nil {
		log.Println(err)
		helper.Logger(helper.LoggerLevelError, "error migrate", err)
		return
	}

	helper.Logger(helper.LoggerLevelInfo, "Running migrate", err)

	if err = m.Up(); err != nil {
		helper.Logger(helper.LoggerLevelError, "Migrate Up Error!!!", err)
		return
	}
	helper.Logger(helper.LoggerLevelInfo, "Migrate Up Done!!!", err)

	return
}
