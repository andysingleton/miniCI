package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Pipeline struct {
	Name             string `json:"name"`
	Filename         string
	ExecutorBackend  string `json:"executors"`
	Dockerfile       string `json:"dockerfile"`
	MiniciBinaryPath string `json:"agentBinaryPath"`
	WebPort          int    `json:"webPort"`
}

func ReadPipelineManifest(manifestFile string) (Pipeline, error) {
	jsonFile, err := os.Open(manifestFile)
	if err != nil {
		return Pipeline{}, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Pipeline{}, err
	}

	var manifest Pipeline
	err = json.Unmarshal(byteValue, &manifest)
	manifest.Filename = manifestFile
	return manifest, err
}
