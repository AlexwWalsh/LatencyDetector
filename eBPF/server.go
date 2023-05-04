package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	// "strings"

	"github.com/gorilla/mux"
)

type Response struct {
	Node `json:"node"`
}

type ProtocolResponse struct {
	ProtocolNode `json:"node"`
}

type Node struct {
	Id         string `json:"id"`
	Ingressing string `json:"ingressing"`
	Egressing  string `json:"egressing"`
}

type ProtocolNode struct {
	Id    string `json:"id"`
	ICMP  string `json:"icmp"`
	TCP   string `json:"tcp"`
	UDP   string `json:"udp"`
	OTHER string `json:"other"`
}

type Address struct {
	IpAddress string `json:"IpAddress"`
}

func main() {
	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/server-status", serverCheck).Methods("GET")
	router.HandleFunc("/countPackets", countPackets).Methods("GET")
	router.HandleFunc("/protocols", protocols).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests

	// Print a message to indicate that the server is listening
	log.Println("Server listening on port 8080")

	data := Address{IpAddress: "192.168.1.216"}
	b, err := json.Marshal(data)
	log.Println(b)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	//Change the below IP address to either localhost (if not working with the VMs), or the machine's IP address that is running the Node Server
	//Example alternative:
	//resp, err := http.Post("http://192.168.1.225:3000/server", "application/json", bytes.NewBuffer(b))
	resp, err := http.Post("http://localhost:3000/server", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body
	fmt.Println(string(body))

	//Change this IP address if running on the VM to that machine's IPv4 address
	//Example alternative:
	//http.ListenAndServe("168.192.18.1:8080", router)
	http.ListenAndServe(":8080", router)
	log.Println("Server listening on port 8080")

}

func serverCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Successfully entered '/server-status' endpointt")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func countPackets(w http.ResponseWriter, r *http.Request) {
	log.Println("Successfully entered '/countPackets' endpoint")
	var response Response
	nodes := prepareResponseCount()

	response.Node = nodes

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}

func protocols(w http.ResponseWriter, r *http.Request) {
	log.Println("Successfully entered '/protocols' endpoint")
	var response ProtocolResponse
	nodes := prepareResponseProtocols()

	response.ProtocolNode = nodes

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}

func runPackageCount() (string, error) {
	cmd := exec.Command("sudo", "go", "run", "./cgroup_skb")

	var stdout bytes.Buffer

	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "Error running packages: ", err
	}

	output := stdout.String()
	return output, nil
}

func runPackageProtocols() (string, error) {
	cmd := exec.Command("sudo", "go", "run", "./xdp")

	var stdout bytes.Buffer

	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "Error running packages: ", err
	}

	output := stdout.String()
	return output, nil
}

func prepareResponseCount() Node {
	// var nodes []Node

	output, err := runPackageCount()

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

	// var i int
	// i = 1
	var node Node
	node.Id = "1"
	node.Ingressing = ingressCount
	node.Egressing = egressCount

	return node
}

func prepareResponseProtocols() ProtocolNode {
	// var nodes []Node

	output, err := runPackageProtocols()

	if err != nil {
		fmt.Println("Error:", err)
	}
	//Ingress is always the first number
	//Egress is always the second number
	//Timing is set to finish counting loop in the same second, but wait a extra second to output egressing so we can tell

	var outputString = string(output)
	splitString := strings.Split(outputString, ", ")

	ICMP := splitString[0]
	TCP := splitString[1]
	UDP := splitString[2]
	OTHER := splitString[3]
	fmt.Println("ICMP:", ICMP)
	fmt.Println("TCP:", TCP)
	fmt.Println("UDP:", UDP)
	fmt.Println("OTHER:", OTHER)

	// var i int
	// i = 1
	var protocolNode ProtocolNode
	protocolNode.Id = "1"
	protocolNode.ICMP = ICMP
	protocolNode.TCP = TCP
	protocolNode.UDP = UDP
	protocolNode.OTHER = OTHER

	return protocolNode
}
