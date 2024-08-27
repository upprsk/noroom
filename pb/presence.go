package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func makeApiNoroomPresence(app *pocketbase.PocketBase, validate *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		type bodyModel struct {
			ClassId     string `json:"class" validate:"required"`
			Fingerprint string `json:"fingerprint" validate:"required"`
			Position    struct {
				Latitude  float64 `json:"latitude" validate:"required"`
				Longitude float64 `json:"longitude" validate:"required"`
			} `json:"position" validate:"required"`
		}

		var body bodyModel
		if err := c.Bind(&body); err != nil {
			return err
		}

		if err := validate.Struct(body); err != nil {
			return err
		}

		fmt.Println("got body:", body)

		klass, err := app.Dao().FindRecordById("classes", body.ClassId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return apis.NewNotFoundError("class not found", body.ClassId)
			}

			return err
		}

		dev, err := app.Dao().FindFirstRecordByData("endDevices", "fingerprint", body.Fingerprint)
		if err != nil {
			return err
		}

		info := apis.RequestInfo(c)
		if !slices.Contains(dev.GetStringSlice("owners"), info.AuthRecord.Id) {
			return apis.NewNotFoundError("The current user is not an owner of the associated device", nil)
		}

		if !klass.GetBool("live") {
			return apis.NewNotFoundError("the given class is not live", body.ClassId)
		}

		latitude := klass.GetFloat("latitude")
		longitude := klass.GetFloat("longitude")
		radius := klass.GetFloat("radius")

		dist := calcDist(latitude, longitude, body.Position.Latitude, body.Position.Longitude)
		if dist > radius {
			return apis.NewNotFoundError("class is to far away", map[string]float64{
				"radius":   radius,
				"distance": dist,
			})
		}

		presenceCollection, err := app.Dao().FindCollectionByNameOrId("classPresenceEntries")
		if err != nil {
			return err
		}

		if err := app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
			pre := models.NewRecord(presenceCollection)
			form := forms.NewRecordUpsert(app, pre)
			form.SetDao(txDao)

			form.LoadData(map[string]any{
				"fingerprint": body.Fingerprint,
				"user":        info.AuthRecord.Id,
				"class":       klass.Id,
			})

			return form.Submit()
		}); err != nil {
			v, ok := err.(validation.Errors)
			if !ok {
				return err
			}

			fingerprintErr, ok := v["fingerprint"]
			if !ok {
				return err
			}

			fingerprint, ok := fingerprintErr.(validation.ErrorObject)
			if !ok {
				return err
			}

			if fingerprint.Code() == "validation_not_unique" {
				return apis.NewBadRequestError(
					"device was already used for presence before for the current class",
					map[string]any{
						"fingerprint": body.Fingerprint,
						"class":       klass.Id,
					},
				)
			}

			return err
		}

		return c.JSON(http.StatusOK, map[string]any{"user": info.AuthRecord.Id, "dist": dist})
	}
}

// https://www.movable-type.co.uk/scripts/latlong.html
//
// This uses the ‘haversine’ formula to calculate the great-circle distance
// between two points – that is, the shortest distance over the earth’s surface
// – giving an ‘as-the-crow-flies’ distance between the points (ignoring any
// hills they fly over, of course!).
// Haversine
// formula:
//
//	a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
//	c = 2 ⋅ atan2( √a, √(1−a) )
//	d = R ⋅ c
//
// where:
//
//	φ is latitude, λ is longitude, R is earth’s radius (mean radius = 6,371km);
//
// > note that angles need to be in radians to pass to trig functions!
//
// JavaScript:
//
//	const R = 6371e3; // metres
//	const φ1 = lat1 * Math.PI/180; // φ, λ in radians
//	const φ2 = lat2 * Math.PI/180;
//	const Δφ = (lat2-lat1) * Math.PI/180;
//	const Δλ = (lon2-lon1) * Math.PI/180;
//
//	const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
//	          Math.cos(φ1) * Math.cos(φ2) *
//	          Math.sin(Δλ/2) * Math.sin(Δλ/2);
//
//	const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
//
//	const d = R * c; // in metres
func calcDist(lat1, lon1, lat2, lon2 float64) float64 {

	const R = 6371e3
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	dphi := (lat2 - lat1) * math.Pi / 180
	dlam := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dphi/2)*math.Sin(dphi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(dlam/2)*math.Sin(dlam/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
