package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ClearBlockchain/onboarding-cli/pkg/gcp"
)

func WriteCredsToEnv(credentials *gcp.Credentials) error {
	gitRepoDir, err := FindGitRepoDir()
	if err != nil {
		log.Fatalf("Failed to find git repo: %v", err)
		return err
	}

	filePath := fmt.Sprintf("%s/.env", gitRepoDir)

	// if .env does not exist, create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create .env file: %v", err)
			return err
		}
		defer file.Close()
	}

	// create creds string
	credsString := fmt.Sprintf(
		"GLIDE_REDIRECT_URI=%s\nGLIDE_CLIENT_ID=%s\nGLIDE_CLIENT_SECRET=%s\n",
		credentials.RedirectURI,
		credentials.ClientID,
		credentials.ClientSecret,
	)

	// append the credentials to the .env file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open .env file: %v", err)
		return nil
	}
	defer file.Close()

	if _, err := file.WriteString(credsString); err != nil {
		log.Fatalf("Failed to write to .env file: %v", err)
		return err
	}

	return nil
}

func FindGitRepoDir() (string, error) {
	// check if the current directory is a git repo
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		log.Errorf("Failed to find git repo: %v", err)
		return "", err
	}

	// remove the newline character
	filePath := strings.Trim(string(out), "\n")
	return filePath, nil
}
