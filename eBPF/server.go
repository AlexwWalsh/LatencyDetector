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

type Node struct {
	Id         string `json:"id"`
	Ingressing string `json:"ingressing"`
	Egressing  string `json:"egressing"`
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
	router.HandleFunc("/data", Data).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests

	// Print a message to indicate that the server is listening
	log.Println("Server listening on port 8080")

	// ifaces, err := net.Interfaces()
	// if err != nil {
	// 	fmt.Println("Error getting addresses: ", err)
	// 	return
	// }

	// for _, i := range ifaces {
	// 	addrs, err := i.Addrs()
	// 	if err != nil {
	// 		fmt.Println("Error getting addresses: ", err)
	// 		return
	// 	}

	// 	for _, addr := range addrs {
	// 		var ip net.IP
	// 		switch v := addr.(type) {
	// 		case *net.IPNet:
	// 			ip = v.IP
	// 			fmt.Println("IP Address: ", ip)
	// 		case *net.IPAddr:
	// 			ip = v.IP
	// 			fmt.Println("IP Address: ", ip)
	// 		}
	// 		fmt.Println("IP Address: ", ip)
	// 	}
	// }

	data := Address{IpAddress: "192.168.1.216"}
	b, err := json.Marshal(data)
	log.Println(b)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}
	resp, err := http.Post("http://192.168.1.225:3000/server", "application/json", bytes.NewBuffer(b))
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

	// Print the response body
	fmt.Println(string(body))
	http.ListenAndServe("192.168.1.216:8080", router)

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

	response.Node = nodes

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

func prepareResponse() Node {
	// var nodes []Node

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

	// var i int
	// i = 1
	var node Node
	node.Id = "1"
	node.Ingressing = ingressCount
	node.Egressing = egressCount

	// node.Id = "3"
	// node.Ingressing = "123"
	// node.Egressing = "6543"
	// nodes = append(nodes, node)

	// node.Id = 3
	// node.Ingressing = "4567"
	// node.Egressing = "398"
	// nodes = append(nodes, node)
	return node
}
