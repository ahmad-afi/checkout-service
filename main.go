package main

import (
	"checkout-service/internal/delivery/http"
	"checkout-service/internal/infrastructure/container"
)

func main() {
	cont := container.NewContainer()
	http.HTTPRouteInit(cont)
}
