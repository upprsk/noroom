package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func main() {
	app := pocketbase.New()

	validate := validator.New(validator.WithRequiredStructEnabled())

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), true))

		e.Router.POST("/api/noroom/tracking", makeApiNoroomTracking(app, validate))

		return nil
	})

	app.OnRecordBeforeCreateRequest("classes").Add(func(e *core.RecordCreateEvent) error {
		admin, _ := e.HttpContext.Get(apis.ContextAdminKey).(*models.Admin)
		if admin != nil {
			return nil // ignore for admins
		}

		info := apis.RequestInfo(e.HttpContext)
		e.Record.Set("owner", info.AuthRecord.Id)
		fmt.Println("owner:", info.AuthRecord.Id)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func makeApiNoroomTracking(app *pocketbase.PocketBase, validate *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		type bodyModel struct {
			UserId       string        `json:"userid"`
			Fingerprint  string        `json:"fingerprint" validate:"required"`
			DeviceData   any           `json:"deviceData"`
			LocationData *LocationInfo `json:"locationData"`
		}

		var body bodyModel
		if err := c.Bind(&body); err != nil {
			return err
		}

		if err := validate.Struct(body); err != nil {
			return err
		}

		dev, err := app.Dao().FindFirstRecordByData("endDevices", "fingerprint", body.Fingerprint)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}

			endDevicesCollection, err := app.Dao().FindCollectionByNameOrId("endDevices")
			if err != nil {
				return err
			}

			dev = models.NewRecord(endDevicesCollection)
		}

		if body.UserId != "" {
			usr, err := app.Dao().FindRecordById("users", body.UserId)
			if err != nil {
				return err
			}

			devOwners := dev.GetStringSlice("owners")
			if len(devOwners) == 0 {
				form := forms.NewRecordUpsert(app, dev)
				form.LoadData(map[string]any{
					"fingerprint":  body.Fingerprint,
					"owners":       []string{usr.Id},
					"deviceData":   body.DeviceData,
					"locationData": body.LocationData,
				})

				if err := form.Submit(); err != nil {
					return err
				}

				return c.JSON(http.StatusOK, map[string]string{
					"device": "new device",
				})
			}

			if !slices.Contains(devOwners, usr.Id) {
				form := forms.NewRecordUpsert(app, dev)
				form.LoadData(map[string]any{
					"owners":       append(devOwners, usr.Id),
					"locationData": body.LocationData,
				})

				if err := form.Submit(); err != nil {
					return err
				}

				return c.JSON(http.StatusOK, map[string]string{
					"device": "updated device",
				})
			}

			form := forms.NewRecordUpsert(app, dev)
			form.LoadData(map[string]any{
				"locationData": body.LocationData,
			})

			if err := form.Submit(); err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"device": "existing device",
			})
		}

		return c.NoContent(http.StatusOK)
	}
}

type LocationInfo struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}
