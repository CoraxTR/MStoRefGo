package main

import (
	"mstorefgo/internal/config"
	moyskladapi "mstorefgo/internal/moyskladapi"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(".env not ready")
	}
	MSRateLimiter := moyskladapi.NewRatelimiter(cfg.RequestCap, cfg.TimeSpan)
	msprocessor := moyskladapi.NewMoySkladProcessor(MSRateLimiter, &cfg.Moyskladapiconfig)

	msprocessor.GetDeliverableOrders()

}
