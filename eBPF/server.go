package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

type Response struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Id         int    `json:"id"`
	Ingressing string `json:"ingressing"`
	Egressing  string `json:"egressing"`
}

func main() {
	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/server-status", serverCheck).Methods("GET")
	router.HandleFunc("/data", Data).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)
}

func serverCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Successfully entered '/server-status' endpointt")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func Data(w http.ResponseWriter, r *http.Request) {
	log.Println("Successfully entered '/data' endpoint")
	var response Response
	nodes := prepareResponse()

	response.Nodes = nodes

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}

func runPackage() (string, error) {
	cmd := exec.Command("sudo", "go", "run", "./cgroup_skb")

	var stdout bytes.Buffer

	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "Error running packages: ", err
	}

	output := stdout.String()
	return output, nil
}

func prepareResponse() []Node {
	var nodes []Node

	output, err := runPackage()

	if err != nil {
		fmt.Println("Error:", err)
	}
	//Ingress is always the first number
	//Egress is always the second number
	//Timing is set to finish counting loop in the same second, but wait a extra second to output egressing so we can tell

	var outputString = string(output)
	splitString := strings.Split(outputString, " ,\n")

	ingressCount := splitString[0]
	egressCount := splitString[1]
	fmt.Println("Ingress:", ingressCount)
	fmt.Println("Egress:", egressCount)

	var i int
	i = 1
	var node Node
	node.Id = i
	node.Ingressing = ingressCount
	node.Egressing = egressCount
	nodes = append(nodes, node)

	// node.Id = 2
	// node.Ingressing = "123"
	// node.Egressing = "6543"
	// nodes = append(nodes, node)

	// node.Id = 3
	// node.Ingressing = "4567"
	// node.Egressing = "398"
	// nodes = append(nodes, node)
	return nodes
}
