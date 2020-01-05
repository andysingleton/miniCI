package main

import (
	"testing"
)

// WorkflowManager test fixture
var getWorkflowsCalled = false

type StubWorkflowManager struct{}

func (StubWorkflowManager) getWorkflows() { getWorkflowsCalled = true }
func (StubWorkflowManager) getAvailableWorkflow(AgentManagerInterface, NetworkManagerInterface) (Workflow, error) {
	return Workflow{
		Name: "foobar",
	}, nil
}

func TestGetWorkflows(t *testing.T) {
	manager := WorkflowManager{}
	manager.ReadWorkflows("test-manifest.json")

	workflowSuccess := false
	for workflow := range manager.Workflows {
		if manager.Workflows[workflow].Name == "Building" {
			workflowSuccess = true
		}
	}
	if workflowSuccess == false {
		t.Errorf("Workflow structure did not contain a test workflow")
	}
}

//func (manager *WorkflowManager) GetAvailableWorkflow(agents []AgentState) (Workflow, error) {
//	manager.UpdateCompletions(agents)
//	for workflow := range manager.Workflows {
//		selectWorkflow := manager.IsWorkflowAvailable(manager.Workflows[workflow])
//		if selectWorkflow == true {
//			return manager.Workflows[workflow], nil
//		}
//	}
//	return Workflow{}, errors.New("No available workflows")
//}

func TestGetAvailableWorkflow_success(t *testing.T) {
	agentHandler := StubAgentManager{}
	agents := agentHandler.GetStates()
	manager := WorkflowManager{}
	err := manager.ReadWorkflows("test-manifest.json")
	check(err)

	workflow, _ := manager.GetAvailableWorkflow(agents)
	if workflow.Name != "Building" {
		t.Errorf("Available workflow was not returned")
	}
}

func TestGetAvailableWorkflow_fail_availability(t *testing.T) {
	agentHandler := StubAgentManagerBusy{}
	agents := agentHandler.GetStates()
	manager := WorkflowManager{}
	err := manager.ReadWorkflows("test-manifest.json")
	check(err)

	workflow, err := manager.GetAvailableWorkflow(agents)
	if err == nil {
		t.Errorf("No workflows should have been returned %s", workflow)
	}
}

func TestGetAvailableWorkflow_fail_no_artefacts(t *testing.T) {
	agentHandler := StubAgentManagerNoArtefacts{}
	agents := agentHandler.GetStates()
	manager := WorkflowManager{}
	err := manager.ReadWorkflows("test-manifest.json")
	check(err)

	workflow, err := manager.GetAvailableWorkflow(agents)
	if err == nil {
		t.Errorf("No workflows should have been returned %s", workflow)
	}
}
