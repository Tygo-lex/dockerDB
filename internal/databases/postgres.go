package databases

import (
	"fmt"
	"os/exec"
)

// PostgresConfig holds configuration for a PostgreSQL container
type PostgresConfig struct {
	Name     string
	Image    string
	Port     string
	User     string
	Password string
	Database string
	Volume   string
    Network  string
}

// SetupPostgresContainer creates and starts a PostgreSQL container
func SetupPostgresContainer(config PostgresConfig) error {
    if err := PullImageWithCLI(config.Image); err != nil {
        return fmt.Errorf("failed to ensure PostgreSQL image: %w", err)
    }
    
    if err := CreateNetworkWithCLI(config.Network); err != nil {
        return fmt.Errorf("failed to create network: %w", err)
    }
    
    args := []string{
        "run", "-d",
        "--name", config.Name,
        "-p", config.Port + ":5432",
        "-v", config.Volume + ":/var/lib/postgresql/data",
        "-e", "POSTGRES_PASSWORD=" + config.Password,
    }

    if config.User != "" {
        args = append(args, "-e", "POSTGRES_USER="+config.User)
    }

    if config.Database != "" {
        args = append(args, "-e", "POSTGRES_DB="+config.Database)
    }

    if config.Network != "" {
        args = append(args, "--network", config.Network)
    }
    
    args = append(args, config.Image)

    cmd := exec.Command("docker", args...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to create PostgreSQL container: %v, output: %s", err, output)
    }

    return nil
}
