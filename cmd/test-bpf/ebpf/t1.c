 /*
  * eBPF program sample.
  * Accepts pointer to first segment packet data as an input parameter.
  * analog of tcpdump -s 1 -d 'dst 1.2.3.4 && udp && dst port 5000‘
  */
#include <stdint.h>
#include <net/ethernet.h>
#include <netinet/ip.h>
#include <netinet/udp.h>
uint64_t entry(void *pkt)
{
	struct ether_header *ether_header = (void *)pkt;
	if (ether_header->ether_type != __builtin_bswap16(0x0800))
		return 0;
	struct iphdr *iphdr = (void *)(ether_header + 1);
	if (iphdr->protocol != 17 || (iphdr->frag_off & 0x1ffff) != 0 ||
	    iphdr->daddr != __builtin_bswap32(0x1020304))
		return 0;
	int hlen = iphdr->ihl * 4;
	struct udphdr *udphdr = (void *)iphdr + hlen;
	if (udphdr->dest != __builtin_bswap16(5000))
		return 0;
	return 1;
}
