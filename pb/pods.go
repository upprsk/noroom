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
	"golang.org/x/net/websocket"
)

const (
	defaultStartTimeout  = time.Second * 20
	defaultDeleteTimeout = defaultStartTimeout
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

		timeout := defaultStartTimeout
		if q := c.QueryParam("timeout"); q != "" {
			timeout, err = time.ParseDuration(q)
			if err != nil {
				apis.NewBadRequestError("invalid timeout", err)
			}
		}

		podId := pod.GetString("podId")
		if err := pm.StartPodById(podId, timeout); err != nil {
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

		timeout := defaultStartTimeout
		if q := c.QueryParam("timeout"); q != "" {
			timeout, err = time.ParseDuration(q)
			if err != nil {
				apis.NewBadRequestError("invalid timeout", err)
			}
		}

		podId := pod.GetString("podId")
		if err := pm.StopPodById(podId, timeout); err != nil {
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

		timeout := defaultStartTimeout
		if q := c.QueryParam("timeout"); q != "" {
			timeout, err = time.ParseDuration(q)
			if err != nil {
				apis.NewBadRequestError("invalid timeout", err)
			}
		}

		podId := pod.GetString("podId")
		if err := pm.KillPodById(podId, timeout); err != nil {
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

func makeApiNoroomPodAttach(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(c echo.Context) error {
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
		stream, err := pm.AttachPodById(podId)
		if err != nil {
			return err
		}

		l := app.Logger()

		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			defer stream.Close()

			go func() {
				defer ws.Close()
				defer stream.Close()

				for {
					buf := make([]byte, 512)
					n, err := stream.Read(buf)
					if err != nil {
						l.Error("error reading from pod stream", "reason", err, "podId", podId)

						// this might be because the container exited
						getAndUpdatePodInspectDataLater(app, pm, id)
						return
					}

					if err := websocket.Message.Send(ws, buf[:n]); err != nil {
						l.Error("error writing to websocket", "reason", err, "podId", podId)
						return
					}
				}
			}()

			for {
				// Read
				var msg []byte
				if err := websocket.Message.Receive(ws, &msg); err != nil {
					app.Logger().Error("error reading from websocket", "reason", err, "podId", podId)
					return
				}

				if _, err := stream.Write(msg); err != nil {
					app.Logger().Error("error writing to pod stream", "reason", err, "podId", podId)
					return
				}
			}
		}).ServeHTTP(c.Response(), c.Request())

		return nil
	}
}

func getAndUpdatePodInspectDataLater(app *pocketbase.PocketBase, pm *pods.PodServerManager, id string) {
	go func() {
		<-time.After(time.Millisecond * 500)

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
