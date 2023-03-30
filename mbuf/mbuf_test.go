package mbuf

import (
	"testing"

	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
)

func TestPktMbufFreeBulk(t *testing.T) {
	eal.InitOnceSafe("test-mbuf", 1)

	mp, err := mempool.CreateMbufPool("test-pool", 1000, 1500)
	if err != nil {
		t.Fatalf("create mempool of mbufs: %v", err)
	}
	defer mp.Free()

	if mp.InUseCount() != 0 {
		t.Fatal("in-use count is not zero")
	}

	mbufs := make([]*Mbuf, 100)
	if err := PktMbufAllocBulk(mp, mbufs); err != nil {
		t.Fatalf("allocate mbufs: %v", err)
	}

	if mp.InUseCount() != 100 {
		t.Fatal("in-use count is not 100")
	}

	PktMbufFreeBulk(mbufs)

	if mp.InUseCount() != 0 {
		t.Fatal("in-use count is not zero")
	}
}
