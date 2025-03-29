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

// MariaDBConfig holds configuration for a MariaDB container
type MariaDBConfig struct {
	Name         string
	Image        string
	Port         string
	RootPassword string
	DatabaseName string
	User         string
	Password     string
	Volume       string
}

// NewMariaDBConfig returns a default MariaDB configuration
func NewMariaDBConfig() *MariaDBConfig {
	return &MariaDBConfig{
		Name:         "mariadb-db",
		Image:        "mariadb:latest",
		Port:         "3306",
		DatabaseName: "mydb",
		Volume:       "mariadb_data",
	}
}

// SetupMariaDBContainer creates and starts a MariaDB container
func SetupMariaDBContainer(config MariaDBConfig) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.40"))
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	if err := PullImageIfNotExists(ctx, cli, config.Image); err != nil {
		return fmt.Errorf("failed to ensure MariaDB image: %w", err)
	}

	// Environment variables for MariaDB
	env := []string{
		"MARIADB_ROOT_PASSWORD=" + config.RootPassword,
		"MARIADB_DATABASE=" + config.DatabaseName,
	}

	if config.User != "" && config.Password != "" {
		env = append(env, "MARIADB_USER="+config.User)
		env = append(env, "MARIADB_PASSWORD="+config.Password)
	}

	// Prepare container configuration
	containerConfig := &container.Config{
		Image: config.Image,
		Env:   env,
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(config.Port + "/tcp"): {},
		},
	}

	// Host configuration with port mapping and volume
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(config.Port + "/tcp"): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: config.Port,
				},
			},
		},
		Binds: []string{config.Volume + ":/var/lib/mysql"},
	}

	// Create container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, config.Name)
	if err != nil {
		return fmt.Errorf("failed to create MariaDB container: %w", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start MariaDB container: %w", err)
	}

	// Wait for the container to be ready
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for MariaDB container to be ready")
		case <-tick:
			inspect, err := cli.ContainerInspect(ctx, resp.ID)
			if err != nil {
				return fmt.Errorf("failed to inspect MariaDB container: %w", err)
			}
			if inspect.State.Running {
				fmt.Println("MariaDB container is running!")
				return nil
			}
		}
	}
}
