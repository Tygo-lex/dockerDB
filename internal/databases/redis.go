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

// RedisConfig holds configuration for a Redis container
type RedisConfig struct {
	Name     string
	Image    string
	Port     string
	Volume   string
	Password string
	Network  string
}

// NewRedisConfig returns a default Redis configuration
func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Name:   "redis",
		Image:  "redis:latest",
		Port:   "6379",
		Volume: "redis_data",
	}
}

// SetupRedisContainer creates and starts a Redis container
// SetupRedisContainer creates and starts a Redis container
func SetupRedisContainer(config *RedisConfig) error {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.40"))
    if err != nil {
        return fmt.Errorf("failed to create Docker client: %w", err)
    }
    if err := PullImageIfNotExists(ctx, cli, config.Image); err != nil {
        return fmt.Errorf("failed to ensure Redis image: %w", err)
    }

    // Create network if specified
    if config.Network != "" {
        // Check if network exists
        networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
        if err != nil {
            return fmt.Errorf("failed to list networks: %w", err)
        }

        networkExists := false
        for _, network := range networks {
            if network.Name == config.Network {
                networkExists = true
                break
            }
        }

        if !networkExists {
            fmt.Printf("Creating network: %s...\n", config.Network)
            _, err = cli.NetworkCreate(ctx, config.Network, types.NetworkCreate{})
            if err != nil {
                return fmt.Errorf("failed to create network: %w", err)
            }
            fmt.Printf("Successfully created network: %s\n", config.Network)
        } else {
            fmt.Printf("Network %s already exists\n", config.Network)
        }
    }

    // Prepare container configuration
    containerConfig := &container.Config{
        Image: config.Image,
        ExposedPorts: map[nat.Port]struct{}{
            nat.Port(config.Port): {},
        },
    }

    // Add password if provided
    if config.Password != "" {
        containerConfig.Cmd = []string{"redis-server", "--requirepass", config.Password}
    }

    // Host configuration with port mapping and volume
    hostConfig := &container.HostConfig{
        PortBindings: nat.PortMap{
            nat.Port(config.Port): []nat.PortBinding{
                {
                    HostIP:   "0.0.0.0",
                    HostPort: config.Port,
                },
            },
        },
        Binds: []string{config.Volume + ":/data"},
    }

    // Network config
    var networkingConfig *network.NetworkingConfig
    if config.Network != "" {
        networkingConfig = &network.NetworkingConfig{
            EndpointsConfig: map[string]*network.EndpointSettings{
                config.Network: {},
            },
        }
    }

    // Create container
    resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, config.Name)
    if err != nil {
        return fmt.Errorf("failed to create Redis container: %w", err)
    }

    // Start the container
    if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
        return fmt.Errorf("failed to start Redis container: %w", err)
    }

    // Wait for the container to be ready
    timeout := time.After(15 * time.Second)
    tick := time.Tick(1 * time.Second)

    for {
        select {
        case <-timeout:
            return fmt.Errorf("timeout waiting for Redis container to be ready")
        case <-tick:
            inspect, err := cli.ContainerInspect(ctx, resp.ID)
            if err != nil {
                return fmt.Errorf("failed to inspect Redis container: %w", err)
            }
            if inspect.State.Running {
                fmt.Println("Redis container is running!")
                return nil
            }
        }
    }
}