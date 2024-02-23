package main

import (
	"github.com/Yougigun/meepshop_q2/internal/repository"
	"github.com/Yougigun/meepshop_q2/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	accountRepo := repository.NewRepository()
	srv := service.Build(logger, accountRepo)
	srv.Engine.Run()
}
