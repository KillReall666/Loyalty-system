package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/KillReall666/Loyalty-system/internal/authentication"
	"github.com/KillReall666/Loyalty-system/internal/config"
	"github.com/KillReall666/Loyalty-system/internal/handlers/addorder"
	"github.com/KillReall666/Loyalty-system/internal/handlers/authorization"
	"github.com/KillReall666/Loyalty-system/internal/handlers/charge"
	"github.com/KillReall666/Loyalty-system/internal/handlers/getbalance"
	"github.com/KillReall666/Loyalty-system/internal/handlers/getorders"
	"github.com/KillReall666/Loyalty-system/internal/handlers/getwithdraw"
	"github.com/KillReall666/Loyalty-system/internal/handlers/registration"
	"github.com/KillReall666/Loyalty-system/internal/interrogator"
	"github.com/KillReall666/Loyalty-system/internal/logger"
	"github.com/KillReall666/Loyalty-system/internal/loyaltysystemservice"
	"github.com/KillReall666/Loyalty-system/internal/storage/postgres"
	"github.com/KillReall666/Loyalty-system/internal/storage/redis"
)

func main() {
	cfg := config.LoadConfig()
	log, err := logger.InitLogger()
	if err != nil {
		panic("couldn't init logger")
	}
	db, err := postgres.NewDB(cfg.DefaultDBConnStr)
	if err != nil {
		log.LogWarning(err)
	}
	log.LogInfo("database connected")
	
	redisClient := redis.NewRedisClient(cfg.RedisAddress)
	pong, err := redisClient.Ping()
	if err != nil {
		log.LogWarning("redis connection error:", err)
	}
	log.LogInfo("Connection to redis established:", pong)

	JWTMiddleware := authentication.JWTMiddleware{
		RedisClient: redisClient,
		Log:         log,
	}

	app := loyaltysystemservice.NewService(db)

	interrog := interrogator.NewInterrogator(db, log)
	go func() {
		for {
			interrog.OrderStatusWorker()
			time.Sleep(1 * time.Second)
		}
	}()
	r := chi.NewRouter()

	r.Use(log.MyLogger)

	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware.JWTMiddleware())
		r.Post(
			"/api/user/orders",
			addorder.NewPutOrderNumberHandler(app, redisClient, log).AddOrderNumberHandler,
		)
		r.Get(
			"/api/user/orders",
			getorders.NewGetOrdersHandler(app, log, interrog).GetOrdersHandler,
		)
		r.Get("/api/user/balance", getbalance.NewGetBalanceHandler(app, log).GetUserBalanceHandler)
		r.Post("/api/user/balance/withdraw", charge.NewChargeHandler(app, log).ChargeHandler)
		r.Get(
			"/api/user/withdrawals",
			getwithdraw.NewGetWithdrawHandler(app, log).GetWithdrawHandler,
		)
	})

	r.Post(
		"/api/user/register",
		registration.NewRegistrationHandler(app, redisClient, log).RegistrationHandler,
	)
	r.Post(
		"/api/user/login",
		authorization.NewAuthorizationHandler(app, redisClient, log).AuthorizationHandler,
	)

	log.LogInfo("starting server at localhost", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.LogWarning(err)
	}
}
