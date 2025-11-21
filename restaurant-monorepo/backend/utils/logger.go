package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	// For production, use NewProduction()
	// For development, use NewDevelopment()
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}
