package main

import (
	"os"

	"github.com/core-api/internal/controller"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error getting env, not coming through %v", err)
	}
	logrus.Info("Successfully loaded env file")
	ctl := controller.NewController()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	err = ctl.Router.Run(":" + port)
}
