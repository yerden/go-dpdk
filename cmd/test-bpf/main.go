package main

/*
#include <rte_mbuf.h>
*/
import "C"

import (
	"flag"
	"fmt"
	"log"
	"unsafe"

	"github.com/yerden/go-dpdk/bpf"
)

var filename = flag.String("f", "", "Specify path to eBPF object")

func main() {
	flag.Parse()

	prm := &bpf.Prm{
		//Insns: bpf.NewInsns(nil),
		XSyms: bpf.NewXSyms([]bpf.XSym{
			&bpf.XSymFunc{
				Name: "rte_pktmbuf_dump",
				Val:  (*bpf.Func)(C.rte_pktmbuf_dump),
				Args: []bpf.Arg{
					&bpf.ArgRaw{
						Size: unsafe.Sizeof(uintptr(0)),
					},
					&bpf.ArgPtrMbuf{
						BufSize: 0,
					},
					&bpf.ArgRaw{
						Size: unsafe.Sizeof(uint32(0)),
					},
				},
			},
		}),
		ProgArg: &bpf.ArgPtrMbuf{
			BufSize: 1500,
		},
	}
	b, err := bpf.ELFLoad(prm, *filename, ".text")
	if err != nil {
		log.Fatalln("unable to create BPF object:", err)
	}
	defer b.Destroy()
	fmt.Printf("obj = %p\n", b)

	data := make([]byte, 150)
	rc := b.Exec(unsafe.Pointer(&data[0]))

	fmt.Printf("rc = %d\n", rc)
}
