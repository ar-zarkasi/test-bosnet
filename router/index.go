package router

import (
	"app/config"
	"app/src/controller"
	"app/src/repository"
	"app/src/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	service = gin.Default()
	validationEngine *validator.Validate
	routerAPI *gin.Engine
)

func init() {
	db := config.GetActiveDB()
	validationEngine = config.GetValidator()
	routerAPI = gin.Default()

	// populate repository
	settingRepository := repository.NewSettings(db)
	counterRepository := repository.NewCounter(db)
	balanceRepository := repository.NewBalance(db)
	historyRepository := repository.NewHistory(db)

	// populate service
	counterService := services.NewCounterService(counterRepository)
	balanceService := services.NewBalanceService(balanceRepository)
	historyService := services.NewHistoryService(historyRepository, balanceService, counterService)
	settingService := services.NewSettings(settingRepository, *counterService, *historyService)

	// initialize
	settingService.InitDeploy()

	// populate controller handler
	groupApi := routerAPI.Group("/")
	bosRouter(groupApi, *balanceService, *counterService, *historyService, validationEngine)
}

func bosRouter(router *gin.RouterGroup, balance services.BalanceService, counter services.CounterService, history services.HistoryService, validation *validator.Validate) {
	controller := controller.NewBosController(balance, counter, history, validation)
	route := router.Group("/api")
		route.PUT("setor", controller.Setor)
		route.PUT("tarik", controller.Tarik)
		route.PUT("transfer", controller.Transfer)
		route.GET("history/:account", controller.ListHistory)
	return
}

func GetRouter() *gin.Engine {
	return routerAPI
}