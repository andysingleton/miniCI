package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func getQualifiedFilename(filename string) string {
	var qualified_filename string

	if filepath.IsAbs(filename) == true {
		qualified_filename = filename
	} else {
		cwd, err := os.Getwd()
		Check(err)
		qualified_filename = filepath.Join(cwd, filename)
	}

	_, err := os.Stat(qualified_filename)
	Check(err)
	return qualified_filename
}

func BuildContainer(filename string, image_name string, workflowId string) {
	cmdStr := fmt.Sprintf("docker build . -f %s -t \"%s\"", filename, image_name)
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	Check(err)

	fmt.Printf("%s", out)
}

func LaunchContainer(imageName string, workflowName string, manifestFile string, workflowId string, miniciBinary string) {
	ChildWorkflowId := fmt.Sprintf("%s-%d", workflowName, rand.Intn(1000))
	containerName := fmt.Sprintf("miniCIAgent-%s", workflowId)
	manifestPath := getQualifiedFilename(manifestFile)
	strippedManifestName := path.Base(manifestFile)
	workflowPath := getQualifiedFilename(workflowName)

	cmdStr := fmt.Sprintf("/usr/bin/docker run --rm -t "+
		"--name %s "+
		"-v %s:/app/miniCIAgent:ro "+
		"-v %s:/app/%s:ro "+
		"-v %s:/app/%s "+
		"%s /app/miniCIAgent -subprocess=yes -id=%s -manifest=docker-demo.json",
		containerName, miniciBinary, manifestPath, strippedManifestName, workflowPath, workflowName,
		imageName, ChildWorkflowId)

	fmt.Printf("%s: Starting child %s\n", workflowId, ChildWorkflowId)
	fmt.Printf("%s: Command is %s\n", workflowId, cmdStr)
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	Check(err)
	fmt.Printf("%s", out)
}
