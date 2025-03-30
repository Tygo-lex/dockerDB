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
    Network  string
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
    
    // Prepare environment variables for auth if enabled
    var env []string
    if config.Auth {
        env = []string{
            "MONGO_INITDB_ROOT_USERNAME=" + config.User,
            "MONGO_INITDB_ROOT_PASSWORD=" + config.Password,
        }
    }
    
    // Create a container
    containerConfig := &container.Config{
        Image: config.Image,
        Env:   env,
        ExposedPorts: map[nat.Port]struct{}{
            nat.Port(config.Port): {},
        },
    }
    
    hostConfig := &container.HostConfig{
        PortBindings: nat.PortMap{
            nat.Port(config.Port): []nat.PortBinding{
                {
                    HostIP:   "0.0.0.0",
                    HostPort: config.Port,
                },
            },
        },
        Binds: []string{config.Volume + ":/data/db"},
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
    
    resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, config.Name)
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