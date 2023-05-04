// The eBPF program will be attached as an packet filter,
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
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf cgroup_skb.c -- -I../headers

func main() {
	// Allow the current process to lock memory for eBPF resources.
	go countPacketsEgress()
	go countPacketsIngress()

	time.Sleep(23 * time.Second)
}

func countPacketsIngress() uint64 {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs1 := bpfObjects{}
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

	log.Println("Counting Ingressing packets...")

	// Read loop reporting the total amount of times the kernel
	// function was entered, once per second.

	var loops = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var value uint64
	for range loops {

		if err := objs1.PktCount.Lookup(uint32(0), &value); err != nil {
			log.Fatalf("reading map: %v", err)
		}
		// log.Printf("number of packets: %d\n", value)
		// fmt.Println("number of packets: \n", value)
		time.Sleep(2 * time.Second)

	}
	fmt.Println(value, ",")
	return 0
}

func countPacketsEgress() uint64 {
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// Get the first-mounted cgroupv2 path.
	cgroupPath, err := detectCgroupPath()
	if err != nil {
		log.Fatal(err)
	}

	// Link the count_egress_packets program to the cgroup.
	l, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Attach:  ebpf.AttachCGroupInetEgress,
		Program: objs.CountEgressPackets,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("Counting Egressing packets...")

	// Read loop reporting the total amount of times the kernel
	// function was entered, once per second.
	var loops = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var value uint64
	for range loops {

		if err := objs.PktCount.Lookup(uint32(0), &value); err != nil {
			log.Fatalf("reading map: %v", err)
		}
		// log.Printf("number of packets: %d\n", value)
		// fmt.Println("number of packets: \n", value)
		time.Sleep(2 * time.Second)

	}
	time.Sleep(1 * time.Second)
	fmt.Println(value)
	return 0
}

// detectCgroupPath returns the first-found mount point of type cgroup2
// and stores it in the cgroupPath global variable.
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
