package main

import (
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/metno/edr-profile/forecast-timeseries/go-example-service/forecastts"
	"github.com/metno/edr-profile/forecast-timeseries/go-example-service/openapi"
)

func main() {
	basePath := flag.String("base-path", "http://localhost:1323", "Base path to use when serving internal URLs.")
	flag.Parse()

	// Create service instance.
	service := forecastts.NewHandler(*basePath)

	e := echo.New()

	forecastts.RegisterHandlers(e, service)
	e.GET("/api", openapi.ServeEcho)

	e.Logger.Fatal(e.Start(":1323"))
}
