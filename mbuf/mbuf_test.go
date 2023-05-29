package mbuf

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
)

func getSample(n int) []byte {
	d := make([]byte, n)
	rand.Read(d)
	return d
}

func TestPktMbufFreeBulk(t *testing.T) {
	eal.InitOnceSafe("test-mbuf", 1)

	mp, err := mempool.CreateMbufPool("test-pool", 1000, 1500)
	assert.NoError(t, err)
	defer mp.Free()

	assert.Zero(t, mp.InUseCount())

	mbufs := make([]*Mbuf, 100)
	assert.NoError(t, PktMbufAllocBulk(mp, mbufs))

	assert.Equal(t, mp.InUseCount(), 100)

	PktMbufFreeBulk(mbufs)
	assert.Zero(t, mp.InUseCount())

	// test on single mbuf
	m := PktMbufAlloc(mp)
	assert.NotNil(t, m)
	assert.Equal(t, mp.InUseCount(), 1)

	clone := m.PktMbufClone(mp)
	assert.NotNil(t, clone)
	assert.Equal(t, mp.InUseCount(), 2)
	clone.PktMbufFree()
	assert.Equal(t, mp.InUseCount(), 1)

	assert.Equal(t, m.HeadRoomSize(), uint16(128))
	assert.Equal(t, m.TailRoomSize(), uint16(1372)) // 1500 - 128
	assert.Equal(t, m.BufLen(), uint16(1500))
	assert.Zero(t, m.PktLen())

	// private area
	priv := &common.CStruct{}
	m.PrivData(priv)
	assert.Zero(t, priv.Len)
	assert.Zero(t, m.PrivSize())
	assert.NotNil(t, priv.Ptr)

	// packet data
	assert.Zero(t, len(m.Data()))
	sample := getSample(100)
	assert.NoError(t, m.PktMbufAppend(sample))
	assert.Equal(t, m.Data(), sample)

	assert.Equal(t, m.Mempool(), mp)

	// ref
	assert.Equal(t, m.RefCntRead(), uint16(1))
	assert.Equal(t, m.RefCntUpdate(1), uint16(2))
	assert.Equal(t, m.RefCntRead(), uint16(2))
	m.RefCntSet(1)
	assert.Equal(t, m.RefCntRead(), uint16(1))
}
