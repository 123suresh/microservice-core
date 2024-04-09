package main

import (
	"os"

	"github.com/core-api/internal/controller"
	"github.com/core-api/internal/repository"
	"github.com/core-api/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error getting env, not coming through %v", err)
	}
	logrus.Info("Successfully loaded env file")
	repo := repository.NewRepo()
	svc := service.NewService(repo)
	ctl := controller.NewController(svc)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	err = ctl.Router.Run(":" + port)
}
