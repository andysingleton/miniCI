package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	mutex sync.Mutex
	//state AgentState
)

type NetworkManagerInterface interface {
	Get() (string, error)
	AddHandler(AgentStateInterface)
	Listen()
	Webport() int
}

type NetworkManager struct {
	Backend string
	WebPort int
}

func (net NetworkManager) Webport() int {
	return net.WebPort
}

func (net NetworkManager) Get() (string, error) {
	switch net.Backend {
	case "local":
		return "172.0.0.1", nil
	case "docker":
		return "172.17.0.1", nil
	}
	return "", errors.New(fmt.Sprintf("Could not handle backend of %s", net.Backend))
}

func (net NetworkManager) AddHandler(state AgentStateInterface) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(state.GetAgentState())
		Check(err)
		fmt.Fprintf(w, string(output))
	})
}

func (net NetworkManager) Listen() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", net.WebPort), nil))
}

func Listener(networkManager NetworkManagerInterface, stateManager AgentStateInterface) {
	fmt.Println(executionId, ": Starting Listener")
	stateManager.InitState(networkManager)
	networkManager.AddHandler(stateManager)
	networkManager.Listen()
}
