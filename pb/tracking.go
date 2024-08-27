package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func makeApiNoroomTracking(app *pocketbase.PocketBase, validate *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		type bodyModel struct {
			UserId      string `json:"userid"`
			Fingerprint string `json:"fingerprint" validate:"required"`
			DeviceData  any    `json:"deviceData"`
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

		realIp := c.Request().Header.Get("X-Real-IP")
		locationInfo, err := getLocationInfoForIp(realIp)
		if err != nil {
			app.Logger().Error(
				"failed to get location info for client address",
				"fingerprint",
				body.Fingerprint,
				"ip",
				realIp,
			)
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
					"locationData": locationInfo,
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
					"locationData": locationInfo,
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
				"locationData": locationInfo,
			})

			if err := form.Submit(); err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"device": "existing device",
			})
		}

		form := forms.NewRecordUpsert(app, dev)
		form.LoadData(map[string]any{
			"fingerprint":  body.Fingerprint,
			"deviceData":   body.DeviceData,
			"locationData": locationInfo,
		})

		if err := form.Submit(); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"device": "anonymous new device",
		})
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

func getLocationInfoForIp(ip string) (*LocationInfo, error) {
	client := &http.Client{
		Timeout: time.Second,
	}

	res, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	info := new(LocationInfo)
	if err := json.NewDecoder(res.Body).Decode(info); err != nil {
		return nil, err
	}

	return info, nil
}
