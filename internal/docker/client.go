package docker

import (
	"fmt"
	"os/exec"
)

// DockerClient is a struct that holds the Docker client configuration.
type DockerClient struct{}

// NewDockerClient creates a new instance of DockerClient.
func NewDockerClient() *DockerClient {
	return &DockerClient{}
}

// RunContainer runs a Docker container with the specified image and options.
func (dc *DockerClient) RunContainer(image string, options []string) error {
	cmd := exec.Command("docker", append([]string{"run"}, options...)...)
	cmd.Args = append(cmd.Args, image)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run container: %v, output: %s", err, output)
	}

	return nil
}

// StopContainer stops a running Docker container by its name or ID.
func (dc *DockerClient) StopContainer(containerID string) error {
	cmd := exec.Command("docker", "stop", containerID)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop container: %v, output: %s", err, output)
	}

	return nil
}

// RemoveContainer removes a Docker container by its name or ID.
func (dc *DockerClient) RemoveContainer(containerID string) error {
	cmd := exec.Command("docker", "rm", containerID)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove container: %v, output: %s", err, output)
	}

	return nil
}

// PullImage pulls a Docker image from the Docker registry.
func (dc *DockerClient) PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull image: %v, output: %s", err, output)
	}

	return nil
}
