package databases

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// PullImageIfNotExists pulls a Docker image if it doesn't exist locally
// For use with Docker API client
func PullImageIfNotExists(ctx context.Context, cli *client.Client, image string) error {
	// Check if image exists locally
	_, _, err := cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		if client.IsErrNotFound(err) {
			// Image doesn't exist locally, pull it
			fmt.Printf("Pulling image: %s...\n", image)
			reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
			if err != nil {
				return fmt.Errorf("failed to pull image: %w", err)
			}
			defer reader.Close()

			// Wait for pull to complete
			_, err = io.Copy(io.Discard, reader)
			if err != nil {
				return fmt.Errorf("error while pulling image: %w", err)
			}
			fmt.Printf("Successfully pulled image: %s\n", image)
		} else {
			return fmt.Errorf("failed to inspect image: %w", err)
		}
	}
	return nil
}

// PullImageWithCLI pulls a Docker image using the docker CLI
func PullImageWithCLI(image string) error {
	// Check if image exists
	checkCmd := exec.Command("docker", "image", "inspect", image)
	if err := checkCmd.Run(); err != nil {
		// Image doesn't exist, pull it
		fmt.Printf("Pulling image: %s...\n", image)
		pullCmd := exec.Command("docker", "pull", image)
		output, err := pullCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to pull image: %v, output: %s", err, output)
		}
		fmt.Printf("Successfully pulled image: %s\n", image)
	}
	return nil
}
func CreateNetworkWithCLI(name string) error {
    if name == "" {
        return nil // No network specified, skip creation
    }
    
    // Check if network exists
    checkCmd := exec.Command("docker", "network", "inspect", name)
    if err := checkCmd.Run(); err == nil {
        // Network exists
        fmt.Printf("Network %s already exists\n", name)
        return nil
    }

    // Network doesn't exist, create it
    fmt.Printf("Creating network: %s...\n", name)
    createCmd := exec.Command("docker", "network", "create", name)
    output, err := createCmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to create network: %v, output: %s", err, output)
    }
    fmt.Printf("Successfully created network: %s\n", name)
    return nil
}