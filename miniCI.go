package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	memberlist "github.com/hashicorp/memberlist"
	"net/http"
	"os"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}
var executionId = uuid.New()

func launcherLoop(agentManager AgentManagerInterface, networkManager NetworkManagerInterface,
	workflowManager WorkflowManagerInterface, pipelineManifest Pipeline) {
	for true {
		// Have all agents exited?
		agentManager.UpdateAgentStates(networkManager)
		fmt.Println(executionId, ": Got stats", agentManager.GetStates())
		if agentManager.CheckAllAgentsDone(networkManager) {
			fmt.Println(executionId, ": All agents have completed. Exiting")
			os.Exit(0)
		}

		// Start new agents if required
		workflowNumber := workflowManager.GetAvailableWorkflowAmount(agentManager)
		fmt.Println(executionId, ": Workflow numbers", workflowNumber)
		for i := 1; i <= workflowNumber; i++ {
			err := agentManager.Start(pipelineManifest)
			Check(err)
		}

		// print web address
		// output logs
		time.Sleep(2)
		fmt.Println(executionId, ": Debug exit")
		os.Exit(0)
	}
}

func main() {
	manifestFile := flag.String("manifest", "default.json", "Name of the Pipeline manifest to build")
	externalArtefacts := flag.String("artifact", "", "Artefacts provided externally to this execution")
	flag.Parse()

	pipelineManifest, err := ReadPipelineManifest(*manifestFile)
	Check(err)
	workflowManager := &WorkflowManager{[]Workflow{}, []string{}, []string{*externalArtefacts}}
	err = workflowManager.ReadWorkflows(*manifestFile)
	Check(err)

	// Start a gossip cluster
	gossipCluster, err := memberlist.Create(memberlist.DefaultLocalConfig())
	Check(err)

	var agentManager AgentManagerInterface
	// todo: add further backend types here, plus associated structs and methods
	if pipelineManifest.ExecutorBackend == "local" {
		agentManager = &AgentManagerLocal{
			*gossipCluster,
			[]AgentState{},
		}
	} else {
		agentManager = &AgentManagerLocal{
			*gossipCluster,
			[]AgentState{},
		}
	}

	// todo: start the artefact service

	// spawn the status webservice
	networkManager := NetworkManager{pipelineManifest.ExecutorBackend, pipelineManifest.WebPort}
	localStateManager := AgentState{}
	go Listener(networkManager, &localStateManager)
	time.Sleep(1 * time.Second)

	launcherLoop(agentManager, networkManager, workflowManager, pipelineManifest)
}
