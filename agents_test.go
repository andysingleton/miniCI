package main

import (
	"github.com/google/uuid"
	"miniCIAgent/memberlist"
	"reflect"
	"testing"
)

// AgentManagerLocal fixture
type StubAgentManager struct {
	agentStates []AgentState
}

func (StubAgentManager) GetMembers() []*memberlist.Node { return []*memberlist.Node{} }
func (StubAgentManager) GetRemoteState(string) (AgentState, error) {
	return AgentState{State: "Testing"}, nil
}
func (handler StubAgentManager) GetStates() []AgentState {
	return handler.agentStates
}
func (handler *StubAgentManager) UpdateAgentStates(NetworkManagerInterface) {
	uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	var agentStates []AgentState
	agentStates = append(agentStates,
		AgentState{
			Ip:          "",
			ExecutionId: uuid,
			State:       "",
			Building:    "",
			Done:        []string{"Testing"},
			Artefacts:   []string{"testArtefact"},
		},
	)
	handler.agentStates = agentStates
}

// AgentManagerLocal fixture: busy nodes
type StubAgentManagerBusy struct{}

func (StubAgentManagerBusy) GetMembers() []*memberlist.Node { return []*memberlist.Node{} }
func (StubAgentManagerBusy) GetRemoteState(string) (AgentState, error) {
	return AgentState{State: "Testing"}, nil
}
func (StubAgentManagerBusy) GetStates() []AgentState {
	uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	var agentStates []AgentState
	agentStates = append(agentStates,
		AgentState{
			Ip:          "",
			ExecutionId: uuid,
			State:       "",
			Building:    "",
			Done:        []string{"Testing"},
			Artefacts:   []string{"testArtefact"},
		},
		AgentState{
			Ip:          "",
			ExecutionId: uuid,
			State:       "",
			Building:    "",
			Done:        []string{"Building"},
			Artefacts:   []string{"foobar.json"},
		},
	)
	return agentStates
}
func (StubAgentManagerBusy) UpdateAgentStates(NetworkManagerInterface) {}

// AgentManagerLocal fixture: artefact requirements not met
type StubAgentManagerNoArtefacts struct{}

func (StubAgentManagerNoArtefacts) GetMembers() []*memberlist.Node { return []*memberlist.Node{} }
func (StubAgentManagerNoArtefacts) GetRemoteState(string) (AgentState, error) {
	return AgentState{State: "Testing"}, nil
}
func (StubAgentManagerNoArtefacts) GetStates() []AgentState {
	uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	var agentStates []AgentState
	agentStates = append(agentStates,
		AgentState{
			Ip:          "",
			ExecutionId: uuid,
			State:       "",
			Building:    "",
			Done:        []string{"Testing"},
			Artefacts:   []string{},
		},
		AgentState{
			Ip:          "",
			ExecutionId: uuid,
			State:       "",
			Building:    "",
			Done:        []string{"Building"},
			Artefacts:   []string{},
		},
	)

	return agentStates
}
func (StubAgentManagerNoArtefacts) UpdateAgentStates(NetworkManagerInterface) {}

func TestGetStates(t *testing.T) {
	uuid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	testAgentState := []AgentState{AgentState{
		Ip:          "",
		ExecutionId: uuid,
		State:       "Testing",
		Building:    "",
		Done:        nil,
		Artefacts:   nil,
	},
	}
	testHandler := AgentManagerLocal{agentStates: testAgentState}
	result := testHandler.GetStates()
	eq := reflect.DeepEqual(result, testAgentState)
	if eq == false {
		t.Errorf("Object does not match expected state: %s", result)
	}
}
