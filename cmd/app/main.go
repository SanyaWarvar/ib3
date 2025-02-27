package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/SanyaWarvar/ib3/pkg/handler"
	"github.com/SanyaWarvar/ib3/pkg/models"
	"github.com/SanyaWarvar/ib3/pkg/repository"
	"github.com/SanyaWarvar/ib3/pkg/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error while load dotenv: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}

	dbNum, err := strconv.Atoi(os.Getenv("CACHE_DB"))
	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}
	redisOptions := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("CACHE_HOST"), os.Getenv("CACHE_PORT")),
		Password: os.Getenv("CACHE_PASSWORD"),
		DB:       dbNum,
	}

	cacheDb, err := repository.NewRedisDb(&redisOptions)
	if err != nil {
		logrus.Fatalf("Error while create connection to cache: %s", err.Error())
	}

	accessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESSTOKENTTL"))
	if err != nil {
		logrus.Fatalf("Errof while parse accessTokenTTL: %s", err.Error())
	}
	refreshTokenTTL, err := time.ParseDuration(os.Getenv("REFRESHTOKENTTL"))
	if err != nil {
		logrus.Fatalf("Errof while parse refreshTokenTTL: %s", err.Error())
	}
	jwtCfg := repository.NewJwtManagerCfg(accessTokenTTL, refreshTokenTTL, os.Getenv("SIGNINGKEY"), jwt.SigningMethodHS256)

	repos := repository.NewRepository(db, cacheDb, jwtCfg)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(models.Server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error while running server: %s", err.Error())
	}

}
