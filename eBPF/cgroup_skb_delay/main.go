// This program demonstrates attaching an eBPF program to a control group.
// The eBPF program will be attached as an egress filter,
// receiving an `__sk_buff` pointer for each outgoing packet.
// It prints the count of total packets every second.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cilium/ebpf"
	// "github.com/cilium/ebpf/internal"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf cgroup_skb.c -- -I../headers

var ErrKeyNotFound = errors.New("key not found")
var inTimeValues []uint64
var outTimeValues []uint64
var conversionMetric uint64 = 1000000

func main() {
	go countPacketDelayIngress()
	go countPacketDelayEgress()
	time.Sleep(10 * time.Second)
	// fmt.Println("Ingress: ", inTimeValues)
	// fmt.Println("Egress: ", outTimeValues)
	// fmt.Println("Time Difference: ", inTimeValues[1]-outTimeValues[1])
	calculateTimeDifferences(inTimeValues, outTimeValues)

}

func calculateTimeDifferences(inTimeValues []uint64, outTimeValues []uint64) uint64 {
	// Remove any indexes with value 0 from the inTimeValues array
	filteredInTimeValues := make([]uint64, 0, len(inTimeValues))
	for _, val := range inTimeValues {
		if val != 0 {
			filteredInTimeValues = append(filteredInTimeValues, val)
		}
	}

	// Remove any indexes with value 0 from the outTimeValues array
	filteredOutTimeValues := make([]uint64, 0, len(outTimeValues))
	for _, val := range outTimeValues {
		if val != 0 {
			filteredOutTimeValues = append(filteredOutTimeValues, val)
		}
	}

	// fmt.Println("Filtered in Time: ", filteredInTimeValues)
	// fmt.Println("filtered out Time: ", filteredOutTimeValues)

	// Subtract the value in index 0 of filteredInTimeValues from each element in filteredOutTimeValues
	var timeDifferences []uint64
	for _, val := range filteredOutTimeValues {
		timeDifferences = append(timeDifferences, val-filteredInTimeValues[0])
	}

	var timeDifferences1 []uint64
	for _, val := range filteredOutTimeValues {
		timeDifferences1 = append(timeDifferences1, filteredInTimeValues[0]-val)
	}

	var timeDifferences2 []uint64
	for _, val := range filteredOutTimeValues {
		timeDifferences2 = append(timeDifferences2, val-filteredInTimeValues[1])
	}

	smallest := ^uint64(0) // initialize smallest to the maximum value of uint64
	for _, val := range timeDifferences {
		if val < smallest {
			smallest = val
		}
	}
	for _, val := range timeDifferences1 {
		if val < smallest {
			smallest = val
		}
	}
	for _, val := range timeDifferences2 {
		if val < smallest {
			smallest = val
		}
	}

	fmt.Println(smallest)
	return smallest
}

func countPacketDelayIngress() uint64 {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs1 := bpfObjects{} //internal
	if err := loadBpfObjects(&objs1, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs1.Close()

	// Get the first-mounted cgroupv2 path.
	cgroupPath, err := detectCgroupPath()
	if err != nil {
		log.Fatal(err)
	}

	// Link the count_ingress_packets program to the cgroup.
	l, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Attach:  ebpf.AttachCGroupInetIngress,
		Program: objs1.CountIngressPackets,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// log.Println("Counting Ingressing packets...")

	// var loops = []int{1, 2, 3, 4, 5, 6, 7, 8}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	// var value uint64
	// var pktInfo PacketInfo

	for range ticker.C {

		var value2 uint64
		if err := objs1.PktInTime.Lookup(uint32(0), &value2); err != nil {
			log.Fatalf("reading map: %v", err)
		}
		// log.Printf("In Time: %d", value2)
		inTimeValues = append(inTimeValues, value2/conversionMetric)
		// time.Sleep(1 * time.Second)
	}

	// log.Printf(" remote_ip= %d, local_ip= %d, remote_port= %U, local_port= %d \n",
	// 	intToIP(pktInfo.RemoteIP4), intToIP(pktInfo.LocalIP4), pktInfo.RemotePort, pktInfo.LocalPort)

	return 0
}

func countPacketDelayEgress() uint64 {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs2 := bpfObjects{} //internal
	if err := loadBpfObjects(&objs2, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs2.Close()

	// Get the first-mounted cgroupv2 path.
	cgroupPath, err := detectCgroupPath()
	if err != nil {
		log.Fatal(err)
	}

	// Link the count_ingress_packets program to the cgroup.
	l, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Attach:  ebpf.AttachCGroupInetEgress,
		Program: objs2.CountEgressPackets,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// log.Println("Counting Egressing packets...")

	// var loops = []int{1, 2, 3, 4, 5, 6, 7, 8}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	// var value uint64
	// var pktInfo PacketInfo

	for range ticker.C {

		var value2 uint64
		if err := objs2.PktOutTime.Lookup(uint32(0), &value2); err != nil {
			log.Fatalf("reading map: %v", err)
		}
		// log.Printf("Out Time: %d", value2)
		outTimeValues = append(outTimeValues, value2/conversionMetric)
		// time.Sleep(1 * time.Second)
	}

	return 0
}

func detectCgroupPath() (string, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// example fields: cgroup2 /sys/fs/cgroup/unified cgroup2 rw,nosuid,nodev,noexec,relatime 0 0
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) >= 3 && fields[2] == "cgroup2" {
			return fields[1], nil
		}
	}

	return "", errors.New("cgroup2 not mounted")
}
