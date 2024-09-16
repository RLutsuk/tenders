package main

import (
	"fmt"
	"log"
	bidDelivery "mymodule/app/internal/bid/delivery"
	bidRep "mymodule/app/internal/bid/repository"
	bidUsecase "mymodule/app/internal/bid/usecase"
	serviceDelivery "mymodule/app/internal/service/delivery"
	tenderDelivery "mymodule/app/internal/tender/delivery"
	tenderRep "mymodule/app/internal/tender/repository"
	tenderUsecase "mymodule/app/internal/tender/usecase"
	userRep "mymodule/app/internal/user/repository"
	"os"
	"regexp"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	postgresConn := os.Getenv("POSTGRES_CONN")
	postgresJDBC := os.Getenv("POSTGRES_JDBC_URL")
	postgresUser := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_DATABASE")

	var config string

	if postgresConn != "" {
		config = postgresConn
	} else if postgresJDBC != "" {
		configJDBC, err := jdbcToGormConfig(postgresJDBC)
		if err == nil {
			config = configJDBC
		}
	} else {
		config = fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s",
			postgresHost, postgresUser, postgresPassword, postgresDB, postgresPort)
	}

	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	userDB := userRep.New(db)

	tenderDB := tenderRep.New(db)
	tenderUC := tenderUsecase.New(tenderDB, userDB)

	bidDB := bidRep.New(db)
	bidUC := bidUsecase.New(bidDB, tenderDB, userDB)

	e := echo.New()
	api := e.Group("/api")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	serviceDelivery.NewDelivery(api)
	tenderDelivery.NewDelivery(api, tenderUC)
	bidDelivery.NewDelivery(api, bidUC)
	e.Logger.Fatal(e.Start(serverAddress))
}

func jdbcToGormConfig(jdbcURL string) (string, error) {
	re := regexp.MustCompile(`jdbc:postgresql://([^:]+):(\d+)/([^?]+)\?user=([^&]+)&password=([^&]+)`)
	matches := re.FindStringSubmatch(jdbcURL)
	if len(matches) != 6 {
		return "", fmt.Errorf("не удалось разобрать строку JDBC")
	}

	host := matches[1]
	port := matches[2]
	dbname := matches[3]
	user := matches[4]
	password := matches[5]

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port), nil
}
