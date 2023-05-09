//go:build ignore
#include "common.h"
#include "vmlinux.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct packet_info {
	__u32 remote_ip4;
	__u32 local_ip4;
	__u16 remote_port;
	__u16 local_port;
};

struct bpf_map_def SEC("maps") pkt_count = {
	.type        = BPF_MAP_TYPE_ARRAY,
	.key_size    = sizeof(u32),
	.value_size  = sizeof(u64),
	.max_entries = 10000,
};

struct bpf_map_def SEC("maps") pkt_in_time = {
	.type        = BPF_MAP_TYPE_HASH,
	.key_size    = sizeof(u64),
	.value_size  = sizeof(u64),
	.max_entries = 10000,
};

struct bpf_map_def SEC("maps") pkt_out_time = {
	.type        = BPF_MAP_TYPE_HASH,
	.key_size    = sizeof(u64),
	.value_size  = sizeof(u64),
	.max_entries = 10000,
};

SEC("cgroup_skb/egress")
int count_egress_packets(struct __sk_buff *skb) {
	u32 key = 0;
	// u32 key2       = 0;
	u64 init_val   = 1;
	u64 time_stamp = bpf_ktime_get_ns();
	u32 sip        = skb->remote_ip4;
	u32 dip        = skb->local_ip4;
	u32 sport      = skb->remote_port;
	u32 dport      = skb->local_port;

	u64 *count = bpf_map_lookup_elem(&pkt_count, &key);
	if (!count) {
		bpf_map_update_elem(&pkt_count, &key, &init_val, BPF_ANY);
		return 1;
	}

	u64 id = ((u64)sip << 32) | ((u64)dip << 0) | ((u64)sport << 16) | ((u64)dport << 0);

	bpf_map_update_elem(&pkt_out_time, &id, &time_stamp, BPF_ANY);

	__sync_fetch_and_add(count, 1);
	return 1;
}

SEC("cgroup_skb/ingress")
int count_ingress_packets(struct __sk_buff *skb) {
	u32 key = 0;
	// u32 key2       = 0;
	u64 init_val   = 1;
	u32 sip        = skb->remote_ip4;
	u32 dip        = skb->local_ip4;
	u32 sport      = skb->remote_port;
	u32 dport      = skb->local_port;
	u64 time_stamp = bpf_ktime_get_ns();

	u64 *count = bpf_map_lookup_elem(&pkt_count, &key);
	if (!count) {
		bpf_map_update_elem(&pkt_count, &key, &init_val, BPF_ANY);
		return 1;
	}

	u64 id = ((u64)sip << 32) | ((u64)dip << 0) | ((u64)sport << 16) | ((u64)dport << 0);

	bpf_map_update_elem(&pkt_in_time, &id, &time_stamp, BPF_ANY);

	__sync_fetch_and_add(count, 1);
	return 1;
}