//go:build ignore

#include "common.h"

#include "bpf_endian.h"
#include "bpf_tracing.h"
#include "vmlinux.h"

#define AF_INET 2

char __license[] SEC("license") = "Dual MIT/GPL";

#define MAX_MAP_ENTRIES 1024

struct bpf_map_def SEC("maps") xdp_stats_map = {
	.type        = BPF_MAP_TYPE_ARRAY,
	.key_size    = sizeof(u32),
	.value_size  = sizeof(u64),
	.max_entries = 4,
};

/*
Attempt to parse the protocol from the packet.
Returns 0 if there is no IPv4 header field; otherwise returns non-zero.
*/
// static __always_inline int parse_ip_src_addr(struct xdp_md *ctx, __u32 *ip_src_addr, __u32 *ip_dest_addr, __u32 *protocol) {
static __always_inline int parse_ip_src_addr(struct xdp_md *ctx, __u32 *protocol) {
	void *data_end = (void *)(long)ctx->data_end;
	void *data     = (void *)(long)ctx->data;

	// First, parse the ethernet header.
	struct ethhdr *eth = data;
	if ((void *)(eth + 1) > data_end) {
		return 0;
	}

	if (eth->h_proto != bpf_htons(ETH_P_IP)) {
		// Check if protocol is ipv4, return if not
		// Remove this check to see packets other than TCP, ICMP, and UDP
		return 0;
	}

	// Then parse the IP header.
	struct iphdr *ip = (void *)(eth + 1);
	if ((void *)(ip + 1) > data_end) {
		return 0;
	}

	*protocol = (__u8)(ip->protocol);
	return 1;
}

SEC("xdp")
int xdp_prog_func(struct xdp_md *ctx) {
	__u32 protocol, index;
	if (!parse_ip_src_addr(ctx, &protocol)) {
		// Not an IPv4 packet, so don't count it.
		goto done;
	}

	if (protocol == 1) { // ICMP
		index = 0;
	} else if (protocol == 6) { // TCP
		index = 1;
	} else if (protocol == 17) { // UDP
		index = 2;
	} else { // Everything else
		index = 3;
	}

	__u32 *pkt_count = bpf_map_lookup_elem(&xdp_stats_map, &index);
	if (!pkt_count) {
		// No entry in the map for this protocol type yet, so set the initial value to 1.
		__u32 init_pkt_count = 1;
		bpf_map_update_elem(&xdp_stats_map, &index, &init_pkt_count, BPF_ANY);
	} else {
		// Entry already exists for this protocol,
		// so increment it atomically using an LLVM built-in.
		__sync_fetch_and_add(pkt_count, 1);
	}

done:
	return XDP_PASS;
}
