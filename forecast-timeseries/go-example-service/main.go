package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

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

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		errorMsg := struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		}{
			Code:        strconv.Itoa(code),
			Description: fmt.Sprintf("%s", err),
		}
		c.JSON(code, &errorMsg)
	}

	e.Logger.Fatal(e.Start(":1323"))
}
