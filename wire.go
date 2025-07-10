//go:build wireinject
// +build wireinject

package main

import (
	"web-page-analyzer/internal/api"
	"web-page-analyzer/internal/persistance"
	"web-page-analyzer/internal/process"
	"web-page-analyzer/internal/service"

	"web-page-analyzer/router"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var StoreSet = wire.NewSet(
	persistance.NewMemoryStore,
	wire.Bind(new(persistance.Store), new(*persistance.InMemoryStore)),
)

func InitializeRouter() *gin.Engine {
	wire.Build(
		StoreSet,
		process.NewPatternExecutor,
		service.NewAnalyzerService,
		api.NewAnalyzerController,
		router.SetupRouter,
	)
	return nil
}
