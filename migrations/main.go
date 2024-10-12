package main

import (
	"checkout-service/internal/helper"
	"checkout-service/internal/infrastructure/container"
	"checkout-service/internal/infrastructure/postgre"
	"checkout-service/internal/utils"
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/color"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var rollback bool
	flag.BoolVar(&rollback, "rollback", false, "")
	var force bool
	flag.BoolVar(&force, "force", false, "")
	var forceMigration bool
	flag.BoolVar(&forceMigration, "forceMigration", false, "")
	var steps = flag.Int("steps", 1, "enter steps")
	var version = flag.Int("version", 0, "enter version")

	flag.Parse()

	err := godotenv.Load(fmt.Sprintf("%s/%s", helper.ProjectRootPath, ".env"))
	if err != nil {
		panic(err)
	}
	container.LoggerInit()
	postgre, err := postgre.Init()
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(postgre.DB, &postgres.Config{
		DatabaseName: utils.EnvString("POSTGRES_DB"),
	})

	if err != nil {
		helper.Logger(helper.LoggerLevelError, "error connect db", err)
		return
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Println(err)
		helper.Logger(helper.LoggerLevelError, "error migrate", err)
		return
	}

	if force {
		if version == nil || *version < 1 {
			helper.Logger(helper.LoggerLevelError, "version required when using force", err)
			return
		}

		fmt.Println("version : ", version)
		if err = m.Force(*version); err != nil {
			helper.Logger(helper.LoggerLevelError, "version required when using force", err)
		}
		helper.Logger(helper.LoggerLevelInfo, "Succeed force version", err)

		version, dirty, err := m.Version()
		if err != nil {
			helper.Logger(helper.LoggerLevelInfo, "failed when check version", err)
		}
		fmt.Println("version :", version)
		fmt.Println("dirty :", dirty)
		return
	}

	helper.Logger(helper.LoggerLevelInfo, "Running migrate", err)
	if rollback {
		if steps == nil {
			color.Println(color.Red("WARNING RUNNING ALL MIGRATION"))

			if !forceMigration {
				color.Println(color.Red("PLEASE USE THIS COMMAND TO FORCE RUNNING ALL MIGRATION"))
				color.Println(color.Yellow("make migrate-rollback steps=0 forceMigration=true"))
			}

			if err = m.Down(); err != nil {
				helper.Logger(helper.LoggerLevelError, "Rollback Error!!!", err)
				return
			}
		} else {

			*steps *= -1
			fmt.Println("step : ", *steps)
			if err = m.Steps(*steps); err != nil {
				helper.Logger(helper.LoggerLevelError, "Rollback Error!!!", err)
				return
			}

		}

		helper.Logger(helper.LoggerLevelInfo, "Rollback Done!!!", err)
	} else {
		if steps == nil {
			color.Println(color.Red("WARNING RUNNING ALL MIGRATION"))

			if !forceMigration {
				color.Println(color.Red("PLEASE USE THIS COMMAND TO FORCE RUNNING ALL MIGRATION"))
				color.Println(color.Yellow("make migrate-up steps=0 forceMigration=true"))
			}

			if err = m.Up(); err != nil {
				helper.Logger(helper.LoggerLevelError, "Migrate Up Error!!!", err)
				return
			}
		} else {
			fmt.Println("step : ", *steps)
			if err = m.Steps(*steps); err != nil {
				helper.Logger(helper.LoggerLevelError, "Migrate Up Error!!!", err)
				return
			}

		}

		helper.Logger(helper.LoggerLevelInfo, "Migrate Up Done!!!", err)
	}
}
