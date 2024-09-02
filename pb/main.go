package main

import (
	"context"
	"errors"
	"log"
	"os"

	"noroom/pb/pods"

	"github.com/go-playground/validator/v10"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func main() {
	app := pocketbase.New()

	podman := pods.NewPodServerManager()

	validate := validator.New(validator.WithRequiredStructEnabled())

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Pre(middlewareLoadAuthContextFromQuery(app))
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), true))

		e.Router.POST("/api/noroom/tracking", makeApiNoroomTracking(app, validate), apis.ActivityLogger(app))

		e.Router.POST("/api/noroom/presence", makeApiNoroomPresence(app, validate), apis.ActivityLogger(app), apis.RequireRecordAuth("users"))

		e.Router.POST("/api/noroom/pod/:id/start", makeApiNoroomPodStart(app, podman), apis.ActivityLogger(app), apis.RequireRecordAuth("users"))
		e.Router.POST("/api/noroom/pod/:id/stop", makeApiNoroomPodStop(app, podman), apis.ActivityLogger(app), apis.RequireRecordAuth("users"))
		e.Router.POST("/api/noroom/pod/:id/kill", makeApiNoroomPodKill(app, podman), apis.ActivityLogger(app), apis.RequireRecordAuth("users"))
		e.Router.POST("/api/noroom/pod/:id/inspect", makeApiNoroomPodInspect(app, podman), apis.ActivityLogger(app), apis.RequireRecordAuth("users"))
		e.Router.GET("/api/noroom/pod/:id/attach", makeApiNoroomPodAttach(app, podman), apis.ActivityLogger(app), apis.RequireRecordAuth("users"))

		if err := checkAndMigrateUsersToHavePods(app); err != nil {
			app.Logger().Error("failed to migrate users", "reason", err)
			return err
		}

		if err := initializePodServerManager(app, podman); err != nil {
			app.Logger().Error("failed to inialize the pod server manager", "reason", err)
		}

		return nil
	})

	app.OnRecordBeforeCreateRequest("classes").Add(makeClassesBeforeCreateRequest())
	app.OnRecordBeforeUpdateRequest("classes").Add(makeClassesBeforeUpdateRequest())

	app.OnRecordBeforeCreateRequest("users").Add(makeUsersBeforeCreateRequest())

	app.OnRecordBeforeCreateRequest("podServers").Add(makePodServersBeforeCreateRequest(podman))
	app.OnRecordBeforeUpdateRequest("podServers").Add(makePodServersBeforeUpdateRequest(podman))
	app.OnRecordAfterDeleteRequest("podServers").Add(makePodServersAfterDeleteRequest(app, podman))

	app.OnRecordBeforeCreateRequest("pods").Add(makePodsBeforeCreateRequest(app, podman))
	app.OnRecordAfterCreateRequest("pods").Add(makePodsAfterCreateRequest(app, podman))
	app.OnRecordBeforeUpdateRequest("pods").Add(makePodsBeforeUpdateRequest(app, podman))
	app.OnRecordAfterDeleteRequest("pods").Add(makePodsAfterDeleteRequest(app, podman))

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func initializePodServerManager(app *pocketbase.PocketBase, pm *pods.PodServerManager) error {
	podServers, err := app.Dao().FindRecordsByExpr("podServers")
	if err != nil {
		return err
	}

	for _, server := range podServers {
		if err := pm.Add(server.Id, server.GetString("address")); err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			app.Logger().Warn("failed to connect to server", "id", server.Id, "name", server.GetString("name"))
			// failed to connect, but keep going (as the server as added to the manager in this case)
		}

		pods, err := app.Dao().FindRecordsByFilter(
			"pods",
			"server={:server}",
			"",
			512,
			0,
			dbx.Params{"server": server.Id},
		)
		if err != nil {
			return err
		}

		for _, pod := range pods {
			if err := pm.AddExistingPodToServerWithoutConnect(server.Id, pod.GetString("podId")); err != nil {
				return err
			}
		}
	}

	return nil
}

// ============================================================================

func checkAndMigrateUsersToHavePods(app *pocketbase.PocketBase) error {
	return app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		users, err := txDao.FindRecordsByExpr("users")
		if err != nil {
			return err
		}

		for _, user := range users {
			if maxPods := user.GetInt("maxPods"); maxPods == 0 {
				user.Set("maxPods", 1)
				if err := txDao.SaveRecord(user); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// ============================================================================

func makeClassesBeforeCreateRequest() func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		admin, _ := e.HttpContext.Get(apis.ContextAdminKey).(*models.Admin)
		if admin != nil {
			return nil // ignore for admins
		}

		info := apis.RequestInfo(e.HttpContext)
		e.Record.Set("owner", info.AuthRecord.Id)

		return nil
	}
}

func makeClassesBeforeUpdateRequest() func(e *core.RecordUpdateEvent) error {
	return func(e *core.RecordUpdateEvent) error {
		admin, _ := e.HttpContext.Get(apis.ContextAdminKey).(*models.Admin)
		if admin != nil {
			return nil // ignore for admins
		}

		original := e.Record.OriginalCopy()

		if original.GetString("owner") != e.Record.GetString("owner") {
			return apis.NewBadRequestError("can't change owner of class", nil)
		}

		return nil
	}
}

// ============================================================================

func makeUsersBeforeCreateRequest() func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		admin, _ := e.HttpContext.Get(apis.ContextAdminKey).(*models.Admin)
		if admin != nil {
			return nil // ignore for admins
		}

		e.Record.Set("role", "student")
		e.Record.Set("maxPods", 1)

		return nil
	}
}

// ============================================================================

func makePodServersBeforeCreateRequest(pm *pods.PodServerManager) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		return pm.Add(e.Record.Id, e.Record.GetString("address"))
	}
}

func makePodServersAfterDeleteRequest(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(e *core.RecordDeleteEvent) error {
	return func(e *core.RecordDeleteEvent) error {
		if err := pm.Del(e.Record.Id); err != nil {
			app.Logger().Error("failed to delete pod server", "podServer", e.Record.Id, "reason", err)
		}

		return nil
	}
}

func makePodServersBeforeUpdateRequest(pm *pods.PodServerManager) func(e *core.RecordUpdateEvent) error {
	return func(e *core.RecordUpdateEvent) error {
		return pm.Update(e.Record.Id, e.Record.GetString("address"))
	}
}

// ============================================================================

func makePodsBeforeCreateRequest(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		info := apis.RequestInfo(e.HttpContext)

		maxPods := info.AuthRecord.GetInt("maxPods")
		pods := info.AuthRecord.GetStringSlice("pods")
		if len(pods)+1 > maxPods {
			return apis.NewForbiddenError("already reached limit of pods for account", map[string]any{
				"maxPods": maxPods,
				"pods":    len(pods),
			})
		}

		serverId := e.Record.GetString("server")
		podName := e.Record.GetString("name")
		podImage := e.Record.GetString("image")

		podId, err := pm.AddNewPodToServer(serverId, podName, podImage)
		if err != nil {
			return err
		}

		e.Record.Set("podId", podId)

		data, err := pm.InspectPodById(podId)
		if err != nil {
			app.Logger().Error("failed to inspect pod after create", "podId", podId, "reason", err)
		} else {
			e.Record.Set("running", data.State.Running)
			e.Record.Set("status", data.State.Status)
		}

		return nil
	}
}

func makePodsAfterCreateRequest(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		info := apis.RequestInfo(e.HttpContext)

		if err := app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
			user, err := txDao.FindRecordById("users", info.AuthRecord.Id)
			if err != nil {
				return err
			}

			pods := user.GetStringSlice("pods")
			pods = append(pods, e.Record.Id)
			user.Set("pods", pods)

			app.Logger().Info(
				"added a pod to user",
				"user",
				info.AuthRecord.Id,
				"pod",
				e.Record.Id,
				"pods",
				pods,
			)

			return txDao.SaveRecord(user)
		}); err != nil {
			return err
		}

		return nil
	}
}

func makePodsAfterDeleteRequest(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(e *core.RecordDeleteEvent) error {
	return func(e *core.RecordDeleteEvent) error {
		serverId := e.Record.GetString("server")
		podId := e.Record.GetString("podId")

		err := pm.DeletePodFromServer(serverId, podId, defaultDeleteTimeout)
		if err != nil {
			app.Logger().Error("failed to delete pod", "podServer", serverId, "pod", e.Record.Id, "reason", err)
		}

		return nil
	}
}

func makePodsBeforeUpdateRequest(app *pocketbase.PocketBase, pm *pods.PodServerManager) func(e *core.RecordUpdateEvent) error {
	return func(e *core.RecordUpdateEvent) error {
		serverId := e.Record.GetString("server")
		podId := e.Record.GetString("podId")

		err := pm.AddExistingPodToServer(serverId, podId)
		if err != nil {
			return err
		}

		data, err := pm.InspectPodById(podId)
		if err != nil {
			app.Logger().Error("failed to inspect pod during update", "podId", podId, "reason", err)
		} else {
			e.Record.Set("running", data.State.Running)
			e.Record.Set("status", data.State.Status)
		}

		return nil
	}
}
