package main

import (
	"fmt"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/metno/edr-profile/forecast-timeseries/go-example-service/metoceants"
	"github.com/metno/edr-profile/forecast-timeseries/go-example-service/openapi"
)

func main() {
	// Create service instance.
	service := metoceants.NewHandler("http://localhost:1323")

	e := echo.New()
	e.Pre(echomiddleware.RemoveTrailingSlash())

	metoceants.RegisterHandlers(e, service)
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
