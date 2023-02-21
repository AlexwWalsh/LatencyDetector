// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/cilium/ebpf"
// 	"github.com/cilium/ebpf/rlimit"
// )

// func main() {

// 	if err := rlimit.RemoveMemlock(); err != nil {
// 		log.Fatalf("Failed to remove rlimit memlock: %s", err)
// 	}

// 	// Create a new eBPF map
// 	packetCounter, err := ebpf.NewMap(&ebpf.MapSpec{
// 		Type:       ebpf.MapType(ebpf.Hash),
// 		KeySize:    4,
// 		ValueSize:  4,
// 		MaxEntries: 100,
// 	})
// 	if err != nil {
// 		fmt.Println("Failed to create map: ", err)
// 		return
// 	}

// 	// Create a new eBPF program
// 	countPackets, err := ebpf.NewProgram(&ebpf.ProgramSpec{
// 		Type:       ebpf.ProgramType(ebpf.Kprobe),
// 		License:    "GPL",
// 		AttachType: ebpf.AttachType(ebpf.Kprobe),
// 		// AttachPoint: "tcp_v4_rcv",
// 	})
// 	if err != nil {
// 		fmt.Println("Failed to create program: ", err)
// 		return
// 	}

// 	// Set the instructions for the program
// 	// err = countPackets.SetInstructions(ebpf.Instructions{
// 	// 	ebpf.LoadAbsolute(ebpf.Reg0, 0),
// 	// 	ebpf.OpAdd(ebpf.Reg0, ebpf.Reg0, 1),
// 	// 	ebpf.StoreMap(ebpf.Reg0, ebpf.Reg0, 0),
// 	// })
// 	if err != nil {
// 		fmt.Println("Failed to set instructions: ", err)
// 		return
// 	}
// 	// countPackets.
// 	// Load the program
// 	err = countPackets.BindMap(packetCounter)
// 	if err != nil {
// 		fmt.Println("Failed to load program: ", err)
// 		return
// 	}
// 	// countPackets.Pin()
// 	// Attach the program to the kprobe
// 	err = countPackets.Pin("tcp_v4_rcv")
// 	if err != nil {
// 		fmt.Println("Failed to attach program: ", err)
// 		return
// 	}
// 	// packetCounter.Iterate()
// 	// Retrieve the count of ingressing packets
// 	packetCount := packetCounter.Lookup(0, 0)
// 	if err != nil {
// 		fmt.Println("Failed to retrieve count: ", err)
// 		return
// 	}

// 	// Print the count of ingressing packets
// 	fmt.Printf("Number of ingressing packets: %d\n", packetCount)
// }
