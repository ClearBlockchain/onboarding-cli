package gcp

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// checkGcloudExists checks if gcloud command is installed and accessible.
func CheckGcloudExists() bool {
	_, err := exec.LookPath("gcloud")
	return err == nil
}

func InstallGcloud() error {
	// Determine the underlying operating system
	switch os := runtime.GOOS; os {
	case "darwin":
		// Install gcloud on macOS
		cmd := exec.Command("brew", "install", "--cask", "google-cloud-sdk")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install gcloud on macOS: %w", err)
		}
	case "linux":
		// Install gcloud on Linux
		cmd := exec.Command("bash", "-c", `curl https://sdk.cloud.google.com | bash`)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install gcloud on Linux: %w", err)
		}
	case "windows":
		// On Windows, we can't install it automatically.
		// We can however direct the user to the installation page
		log.Info("Please download and install Google Cloud SDK from: https://cloud.google.com/sdk/docs/quickstart-windows and try again.")
	default:
		return fmt.Errorf("%s is not supported", os)
	}

	return nil
}

func GetGCPProjects() ([]string, error) {
	// Check if gcloud command exists
	if !CheckGcloudExists() {
		err := fmt.Errorf("gcloud command not found")
		log.Error(err)
		return nil, err
	}

	// Create a context with a timeout that will abort the command
	// if it does not finish in 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// use gcloud to get the list of projects
	cmd := exec.CommandContext(ctx, "gcloud", "projects", "list", "--format=value(projectId)")
	out, err := cmd.Output()
	if err != nil {
		log.Error("Error executing gcloud command: ", err)
		return nil, err
	}

	// If the output is empty, return an empty slice of strings
	if len(out) == 0 {
		log.Warn("No GCP projects found")
		return []string{}, nil
	}

	// split the output into a slice of strings
	projects := strings.Split(string(out), "\n")

	// remove the last empty string
	if len(projects) > 0 && projects[len(projects)-1] == "" {
		projects = projects[:len(projects)-1]
	}

	return projects, nil
}
