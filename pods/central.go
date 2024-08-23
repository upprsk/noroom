package main

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Central struct {
	docker *client.Client
}

func NewCentral(ctx context.Context) (*Central, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Central{
		docker: docker,
	}, nil
}

func (c *Central) Close() {
	c.docker.Close()
}

func (c *Central) ContainersList(ctx context.Context) ([]string, error) {
	containers, err := c.docker.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, container := range containers {
		ids = append(ids, container.ID)
	}

	return ids, nil
}

func (c *Central) ContainerCreate(ctx context.Context) (string, error) {
	resp, err := c.docker.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"sh"},
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
