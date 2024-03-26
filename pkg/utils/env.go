package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func WriteCredsToEnv(credentials map[string]string) error {
	gitRepoDir, err := FindGitRepoDir()
	if err != nil {
		log.Fatalf("Failed to find git repo: %v", err)
		return err
	}

	filePath := fmt.Sprintf("%s/.env", gitRepoDir)

	// append variables to the .env file
	if err := AppendEnvVars(
		filePath,
		credentials,
	); err != nil {
		log.Fatalf("Failed to append env vars: %v", err)
		return err
	}

	return nil
}

func AppendEnvVars(fileName string, vars map[string]string) error {
	// Open the file in read-write mode or create it if it does not exist
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file line by line and store it in a slice
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check if the keys exist in the file and update them
	for key, value := range vars {
		exists := false
		for i, line := range lines {
			if strings.HasPrefix(line, key+"=") {
				lines[i] = key + "=" + value
				exists = true
				break
			}
		}

		// If the key does not exist in the file, append it
		if !exists {
			lines = append(lines, key+"="+value)
		}
	}

	// Write the updated content to the file
	file.Seek(0, 0)
	file.Truncate(0)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
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
