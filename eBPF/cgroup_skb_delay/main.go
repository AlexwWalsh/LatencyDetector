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
	"sort"
	"strings"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf cgroup_skb.c -- -I../headers

var ErrKeyNotFound = errors.New("key not found")

var conversionMetric uint64 = 1000000
var map1 = make(map[uint64]uint64)
var map2 = make(map[uint64]uint64)

func main() {

	go countPacketDelayIngress()
	go countPacketDelayEgress()
	time.Sleep(15 * time.Second)

	resultMap := subtractMaps(map1, map2)

	resultArray := getPerformanceArray(resultMap)

	fmt.Println(resultArray[0], resultArray[1], resultArray[2])
}

func getPerformanceArray(m map[uint64]uint64) [3]uint64 {
	var nums []uint64
	for _, val := range m {
		if val >= 10 && val <= 9999 {
			nums = append(nums, val)
		}
	}
	if len(nums) == 0 {
		return [3]uint64{0, 0, 0}
	}
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	highest := nums[len(nums)-1]
	var sum uint64
	for _, num := range nums {
		sum += num
	}
	average := sum / uint64(len(nums))
	lowest := nums[0]
	return [3]uint64{highest, average, lowest}
}

func subtractMaps(map1, map2 map[uint64]uint64) map[uint64]uint64 {
	resultMap := make(map[uint64]uint64)

	// Iterate over the keys in map1
	for key, value1 := range map1 {
		// Check if the key is also in map2
		if value2, ok := map2[key]; ok {
			// Subtract the value from map2 from the value from map1
			result := value1 - value2
			// Store the result in the output map
			resultMap[key] = result
		}
	}

	return resultMap
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
	var loops = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var s string
	for range loops {
		contents, err := formatMapContentsIngress(objs1.PktInTime)
		if err != nil {
			log.Printf("Error reading map: %s", err)
			continue
		}
		s = contents // assign the loop variable to the outer variable
		time.Sleep(2 * time.Second)
	}
	fmt.Println(s)

	// log.Printf(" remote_ip= %d, local_ip= %d, remote_port= %U, local_port= %d \n",
	// 	intToIP(pktInfo.RemoteIP4), intToIP(pktInfo.LocalIP4), pktInfo.RemotePort, pktInfo.LocalPort)

	return 0
}

func formatMapContentsIngress(m *ebpf.Map) (string, error) {
	var (
		sb  strings.Builder
		key uint64
		val uint64
	)
	iter := m.Iterate()
	for iter.Next(&key, &val) {
		sb.WriteString(fmt.Sprintf("%d=>%d\n", key, val))
		map1[key] = val / conversionMetric
	}
	return sb.String(), iter.Err()
}

func formatMapContentsEgress(m *ebpf.Map) (string, error) {
	var (
		sb   strings.Builder
		key  uint64
		val1 uint64
	)
	iter := m.Iterate()
	for iter.Next(&key, &val1) {
		sb.WriteString(fmt.Sprintf("%d=>%d\n", key, val1))
		map2[key] = val1 / conversionMetric
	}
	return sb.String(), iter.Err()
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

	var loops = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var s string
	for range loops {
		contents, err := formatMapContentsEgress(objs2.PktOutTime)
		if err != nil {
			log.Printf("Error reading map: %s", err)
			continue
		}
		s = contents // assign the loop variable to the outer variable
		time.Sleep(2 * time.Second)
	}
	fmt.Println(s)

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
