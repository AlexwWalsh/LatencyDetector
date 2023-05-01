//go:build ignore
#include "vmlinux.h"
char __license[] SEC("license") = "Dual MIT/GPL";


struct packet_info {
	__u32 hash;
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

struct bpf_map_def SEC("maps") pkt_info_map = {
	.type        = BPF_MAP_TYPE_ARRAY,
	.key_size    = sizeof(u32),
	.value_size  = sizeof(struct packet_info),
	.max_entries = 1000000,
};

SEC("cgroup_skb/egress")
int count_egress_packets(struct __sk_buff *skb) {
	u32 key      = 0;
	u64 init_val = 1;

	u64 *count = bpf_map_lookup_elem(&pkt_count, &key);
	if (!count) {
		bpf_map_update_elem(&pkt_count, &key, &init_val, BPF_ANY);
		return 1;
	}
	__sync_fetch_and_add(count, 1);
	return 1;
}
SEC("cgroup_skb/ingress")
int count_ingress_packets(struct __sk_buff *skb) {
	u32 key      = 0;
	u64 init_val = 1;


	u64 *count = bpf_map_lookup_elem(&pkt_count, &key);
	if (!count) {
		bpf_map_update_elem(&pkt_count, &key, &init_val, BPF_ANY);
		return 1;
	}
	__sync_fetch_and_add(count, 1);

	struct packet_info pkt_info = {
		.hash        = skb->hash,
		.remote_ip4  = skb->remote_ip4,
		.local_ip4   = skb->local_ip4,
		.remote_port = skb->remote_port,
		.local_port  = skb->local_port,
	};
	//&map, &key, &value, &BPF_VALUE 
	bpf_map_update_elem(&pkt_info_map, &pkt_info.hash, &pkt_info, BPF_ANY);
	return 1;
}
