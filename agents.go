package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	memberlist "github.com/hashicorp/memberlist"
	"io/ioutil"
	"os/exec"
)

type AgentStateInterface interface {
	InitState(NetworkManagerInterface)
	SetStatus(string)
	SetBuilding(string)
	AddDone(string)
	AddArtefact(string)
	GetAgentState() AgentState
	SetPendingWorkflow(string)
	PromoteToBuilding(string)
	PromoteToDone(string)
}

type AgentState struct {
	Ip          string
	ExecutionId uuid.UUID
	State       string
	Building    string
	Pending     string
	Done        []string
	Artefacts   []string
}

func (agentState *AgentState) SetPendingWorkflow(workflowName string) {
	agentState.Pending = workflowName
}

func (agentState *AgentState) PromoteToDone(workflowName string) {
	agentState.Done = append(agentState.Done, workflowName)
	agentState.Building = ""
}

func (agentState *AgentState) PromoteToBuilding(workflowName string) {
	agentState.Building = workflowName
	agentState.Pending = ""
}

func (st AgentState) GetAgentState() AgentState {
	return st
}

func (st *AgentState) InitState(ipGetter NetworkManagerInterface) {
	var err error

	mutex.Lock()
	st.Ip, err = ipGetter.Get()
	Check(err)
	st.ExecutionId = executionId
	st.State = "Starting"
	mutex.Unlock()
}

func (st *AgentState) SetStatus(newStatus string) {
	mutex.Lock()
	st.State = newStatus
	mutex.Unlock()
}

func (st *AgentState) SetBuilding(workflowName string) {
	mutex.Lock()
	st.Building = workflowName
	mutex.Unlock()
}

func (st *AgentState) AddDone(workflowName string) {
	mutex.Lock()
	st.Done = append(st.Done, workflowName)
	mutex.Unlock()
}

func (st *AgentState) AddArtefact(artefact string) {
	mutex.Lock()
	st.Artefacts = append(st.Artefacts, artefact)
	mutex.Unlock()
}

func listener(networkManager NetworkManagerInterface, stateManager AgentStateInterface) {
	fmt.Println(executionId, ": Starting Listener")

	stateManager.InitState(networkManager)

	networkManager.AddHandler(stateManager)
	networkManager.Listen()
}

type AgentManagerInterface interface {
	GetRemoteState(string) (AgentState, error)
	GetMembers() []*memberlist.Node
	GetStates() []AgentState
	UpdateAgentStates(NetworkManagerInterface)
	Start(Pipeline) error
	CheckAllAgentsDone(NetworkManagerInterface) bool
}

type AgentManagerLocal struct {
	Gossip      memberlist.Memberlist
	AgentStates []AgentState
}

func (handler AgentManagerLocal) GetStates() []AgentState {
	return handler.AgentStates
}

func (handler AgentManagerLocal) Start(pipeline Pipeline) error {
	fmt.Println(executionId, ": Starting agent")
	args := fmt.Sprintf("--manifest %s", pipeline.Filename)
	cmd := exec.Command(pipeline.MiniciBinaryPath, args)
	err := cmd.Start()
	return err
}

func (handler AgentManagerLocal) CheckAllAgentsDone(networkManager NetworkManagerInterface) bool {
	handler.UpdateAgentStates(networkManager)
	for agent := range handler.AgentStates {
		// If any node is in a working state, we return false
		switch handler.AgentStates[agent].State {
		case "Starting":
			return false
		case "Building":
			return false
		}
	}
	return true
}

func (handler AgentManagerLocal) GetRemoteState(url string) (AgentState, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return AgentState{}, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return AgentState{}, err
	}

	var result AgentState
	err = json.Unmarshal([]byte(body), &result)
	return result, err
}

func (handler AgentManagerLocal) GetMembers() []*memberlist.Node {
	return handler.Gossip.Members()
}

func (handler *AgentManagerLocal) UpdateAgentStates(networkManager NetworkManagerInterface) {
	memberList := handler.GetMembers()
	var agents []AgentState

	for member := range memberList {
		var agentState AgentState
		connectionString := fmt.Sprintf("http://%s:%d", memberList[member], networkManager.Webport())
		agentState, err := handler.GetRemoteState(connectionString)
		Check(err)
		agents = append(agents, agentState)
	}
	handler.AgentStates = agents
}
