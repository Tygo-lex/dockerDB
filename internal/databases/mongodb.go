package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type MongoDBConfig struct {
	Name     string
	Image    string
	Port     string
	Volume   string
	User     string
	Password string
	Auth     bool
}

func NewMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		Image:  "mongo:latest",
		Name:   "mongodb",
		Port:   "27017",
		Volume: "/data/db",
	}
}

func SetupMongoDB(ctx context.Context, config *MongoDBConfig) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.40"))
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	if err := PullImageIfNotExists(ctx, cli, config.Image); err != nil {
		return fmt.Errorf("failed to ensure MongoDB image: %w", err)
	}
	// Create a container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: config.Image,
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(config.Port): {},
		},
	}, &container.HostConfig{
		Binds: []string{config.Volume + ":/data/db"},
	}, nil, nil, config.Name)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB container: %w", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	// Wait for the container to be ready
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for MongoDB container to be ready")
		case <-tick:
			inspect, err := cli.ContainerInspect(ctx, resp.ID)
			if err != nil {
				return fmt.Errorf("failed to inspect MongoDB container: %w", err)
			}
			if inspect.State.Running {
				return nil
			}
		}
	}
}
