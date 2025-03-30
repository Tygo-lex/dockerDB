package databases

import (
	"fmt"
	"os/exec"
)

// MySQLConfig holds configuration for a MySQL container
type MySQLConfig struct {
	Name         string
	Image        string
	Port         string
	RootPassword string
	DatabaseName string
	User         string
	Password     string
	Volume       string
    Network      string
}

// SetupMySQLContainer creates and starts a MySQL container
func SetupMySQLContainer(config MySQLConfig) error {
	if err := PullImageWithCLI(config.Image); err != nil {
		return fmt.Errorf("failed to ensure MySQL image: %w", err)
	}
	if err := CreateNetworkWithCLI(config.Network); err != nil {
        return fmt.Errorf("failed to create network: %w", err)
    }
	args := []string{
		"run", "-d",
		"--name", config.Name,
		"-p", config.Port + ":3306",
		"-v", config.Volume + ":/var/lib/mysql",
		"-e", "MYSQL_ROOT_PASSWORD=" + config.RootPassword,
		"-e", "MYSQL_DATABASE=" + config.DatabaseName,
	}

	if config.User != "" && config.Password != "" {
		args = append(args, "-e", "MYSQL_USER="+config.User)
		args = append(args, "-e", "MYSQL_PASSWORD="+config.Password)
	}

	if config.Network != "" {
        args = append(args, "--network", config.Network)
    }

	args = append(args, config.Image)

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create MySQL container: %v, output: %s", err, output)
	}

	return nil
}
