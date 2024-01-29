package openapi

import (
	"net/http"

	"github.com/labstack/echo/v4"

	_ "embed"
)

//go:embed metocean-ts-bundle.yaml
var schema []byte

// Serve sends the openapi schema to the requestor.
func Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/yaml")
	w.Write(schema)
}

func ServeEcho(ctx echo.Context) error {
	return ctx.Blob(http.StatusOK, "application/yaml", schema)
}
