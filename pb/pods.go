package main

import (
	"net/http"
	"noroom/pb/pods"
	"noroom/rpc"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
)

func makeApiNoroomPodStart(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(c echo.Context) error {
	return func(c echo.Context) error {
		info := apis.RequestInfo(c)

		id := c.PathParam("id")
		if id == "" {
			return apis.NewBadRequestError("missing id", nil)
		}

		pod, err := app.Dao().FindRecordById("pods", id)
		if err != nil {
			return err
		}

		canAccess, err := app.Dao().CanAccessRecord(pod, info, pod.Collection().UpdateRule)
		if !canAccess {
			return apis.NewForbiddenError("", err)
		}

		podId := pod.GetString("podId")
		if err := pm.StartPodById(podId); err != nil {
			return err
		}

		getAndUpdatePodInspectDataLater(app, pm, id)

		return c.NoContent(http.StatusOK)
	}
}

func makeApiNoroomPodStop(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(c echo.Context) error {
	return func(c echo.Context) error {
		info := apis.RequestInfo(c)

		id := c.PathParam("id")
		if id == "" {
			return apis.NewBadRequestError("missing id", nil)
		}

		pod, err := app.Dao().FindRecordById("pods", id)
		if err != nil {
			return err
		}

		canAccess, err := app.Dao().CanAccessRecord(pod, info, pod.Collection().UpdateRule)
		if !canAccess {
			return apis.NewForbiddenError("", err)
		}

		podId := pod.GetString("podId")
		if err := pm.StopPodById(podId); err != nil {
			return err
		}

		getAndUpdatePodInspectDataLater(app, pm, id)

		return c.NoContent(http.StatusOK)
	}
}

func makeApiNoroomPodKill(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(c echo.Context) error {
	return func(c echo.Context) error {
		info := apis.RequestInfo(c)

		id := c.PathParam("id")
		if id == "" {
			return apis.NewBadRequestError("missing id", nil)
		}

		pod, err := app.Dao().FindRecordById("pods", id)
		if err != nil {
			return err
		}

		canAccess, err := app.Dao().CanAccessRecord(pod, info, pod.Collection().UpdateRule)
		if !canAccess {
			return apis.NewForbiddenError("", err)
		}

		podId := pod.GetString("podId")
		if err := pm.KillPodById(podId); err != nil {
			return err
		}

		getAndUpdatePodInspectDataLater(app, pm, id)

		return c.NoContent(http.StatusOK)
	}
}

func makeApiNoroomPodInspect(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(c echo.Context) error {
	return func(c echo.Context) error {
		info := apis.RequestInfo(c)

		id := c.PathParam("id")
		if id == "" {
			return apis.NewBadRequestError("missing id", nil)
		}

		pod, err := app.Dao().FindRecordById("pods", id)
		if err != nil {
			return err
		}

		canAccess, err := app.Dao().CanAccessRecord(pod, info, pod.Collection().ViewRule)
		if !canAccess {
			return apis.NewForbiddenError("", err)
		}

		podId := pod.GetString("podId")
		data, err := pm.InspectPodById(podId)
		if err != nil {
			return err
		}

		if err := updatePodInspectData(app, id, data); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, data)
	}
}

func getAndUpdatePodInspectDataLater(app *pocketbase.PocketBase, pm *pods.PodServerManager, id string) {
	go func() {
		<-time.After(time.Second)

		if err := getAndUpdatePodInspectData(app, pm, id); err != nil {
			app.Logger().Error("failed to get and update pod state", "id", id, "reason", err)
		}
	}()
}

func getAndUpdatePodInspectData(app *pocketbase.PocketBase, pm *pods.PodServerManager, id string) error {
	pod, err := app.Dao().FindRecordById("pods", id)
	if err != nil {
		return err
	}

	podId := pod.GetString("podId")
	data, err := pm.InspectPodById(podId)
	if err != nil {
		return err
	}

	return updatePodInspectData(app, id, data)
}

func updatePodInspectData(
	app *pocketbase.PocketBase,
	id string,
	data *rpc.ContainerInspectResult,
) error {
	return app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		pod, err := txDao.FindRecordById("pods", id)
		if err != nil {
			return err
		}

		form := forms.NewRecordUpsert(app, pod)
		form.SetDao(txDao)

		s := data.State
		form.LoadData(map[string]any{
			"running": s.Running,
			"status":  s.Status,
		})

		return form.Submit()
	})
}
