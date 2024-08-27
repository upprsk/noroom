package hub

import (
	"context"
	"log"
	"noroom/rpc"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Hub struct {
	docker *client.Client
}

func NewHub(ctx context.Context) (*Hub, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Hub{
		docker: docker,
	}, nil
}

func (h *Hub) Create(ctx context.Context, name, image string, cmd []string, env []string) (string, error) {
	log.Printf("Create(name=%v, image=%v, cmd=%v, env=%v)", name, image, cmd, env)

	resp, err := h.docker.ContainerCreate(ctx, &container.Config{
		Cmd:          cmd,
		Env:          env,
		Image:        image,
		WorkingDir:   "/home",
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
	}, nil, nil, nil, name)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (h *Hub) Start(ctx context.Context, id string) error {
	log.Printf("Start(id=%v)", id)

	if err := h.docker.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		log.Println("Start err:", err)
		return err
	}

	return nil
}

func (h *Hub) Stop(ctx context.Context, id string) error {
	log.Printf("Stop(id=%v)", id)

	if err := h.docker.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
		log.Println("Stop err:", err)
		return err
	}

	return nil
}

func (h *Hub) Kill(ctx context.Context, id, signal string) error {
	log.Printf("Kill(id=%v)", id)

	if err := h.docker.ContainerKill(ctx, id, signal); err != nil {
		log.Println("Kill err:", err)
		return err
	}

	return nil
}

func (h *Hub) Delete(ctx context.Context, id string) error {
	log.Printf("Delete(id=%v)", id)

	err := h.docker.ContainerRemove(ctx, id, container.RemoveOptions{})
	if err != nil {
		log.Println("Delete err:", err)
		return err
	}

	return nil
}

func (h *Hub) Inspect(ctx context.Context, id string) (*rpc.ContainerInspectResult, error) {
	log.Printf("Inspect(id=%v)", id)

	data, err := h.docker.ContainerInspect(ctx, id)
	if err != nil {
		log.Println("Inspect err:", err)
		return nil, err
	}

	return &rpc.ContainerInspectResult{
		Id:         id,
		Name:       data.Name,
		Path:       data.Path,
		Args:       data.Args,
		Image:      data.Image,
		Created:    data.Created,
		SizeRw:     data.SizeRw,
		SizeRootFs: data.SizeRootFs,
		State: rpc.ContainerState{
			Status:     data.State.Status,
			Running:    data.State.Running,
			Paused:     data.State.Paused,
			Restarting: data.State.Restarting,
			OOMKilled:  data.State.OOMKilled,
			Dead:       data.State.Dead,
			Pid:        data.State.Pid,
			ExitCode:   data.State.ExitCode,
			Error:      data.State.Error,
			StartedAt:  data.State.StartedAt,
			FinishedAt: data.State.FinishedAt,
		},
	}, nil
}
