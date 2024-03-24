package gcp

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetGCPProjects() []string {
	// use gcloud to get the list of projects
	cmd := exec.Command("gcloud", "projects", "list", "--format=value(projectId)")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// split the output into a slice of strings
	projects := strings.Split(string(out), "\n")

	// remove the last empty string
	projects = projects[:len(projects)-1]

	return projects
}
