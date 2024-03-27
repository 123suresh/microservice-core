package main

import (
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
}
