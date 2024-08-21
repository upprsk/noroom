package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), true))

		e.Router.POST("/api/noroom/tracking", func(c echo.Context) error {
			type bodyModel struct {
				UserId      string `json:"userid"`
				Fingerprint string `json:"fingerprint"`
			}

			var body bodyModel
			if err := c.Bind(&body); err != nil {
				return err
			}

			app.Logger().Debug("got body:", "body", body)

			return c.NoContent(http.StatusOK)
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
