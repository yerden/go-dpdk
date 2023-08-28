package ethdev

import (
	"bytes"
	"syscall"
	"testing"

	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
	"github.com/yerden/go-dpdk/ring"
)

func TestPortFromRing(t *testing.T) {
	eal.InitOnceSafe("test", 4)
	r, err := ring.NewEAL("test_ring", 128)
	assert(t, err == nil, err)
	defer r.Free()

	pid, err := FromRing(r)
	assert(t, err == nil)
	defer pid.Close()

	err = pid.DevConfigure(1, 1)
	assert(t, err == nil)

	mp, err := mempool.CreateMbufPool("test_mp", 1024, 2048)
	assert(t, err == nil, err)
	defer mp.Free()

	err = pid.RxqSetup(0, 128, mp)
	assert(t, err == nil, err)

	err = pid.TxqSetup(0, 128)
	assert(t, err == nil, err)
}

func TestPortFromRings(t *testing.T) {
	eal.InitOnceSafe("test", 4)
	r1, err := ring.NewEAL("test_ring1", 128)
	assert(t, err == nil, err)
	defer r1.Free()
	r2, err := ring.NewEAL("test_ring2", 128)
	assert(t, err == nil, err)
	defer r2.Free()

	pid, err := FromRings("test_port", []*ring.Ring{r1, r2}, []*ring.Ring{r1, r2}, 0)
	assert(t, err == nil)
	defer pid.Close()

	err = pid.DevConfigure(2, 2)
	assert(t, err == nil)

	mp, err := mempool.CreateMbufPool("test_mp", 1024, 2048)
	assert(t, err == nil, err)
	defer mp.Free()

	err = pid.RxqSetup(0, 128, mp)
	assert(t, err == nil, err)

	err = pid.TxqSetup(0, 128)
	assert(t, err == nil, err)
}

func TestMACAddr(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	var a MACAddr
	err := pid.MACAddrGet(&a)
	assert(t, err == nil, err)

	hwAddr := a.HardwareAddr()
	assert(t, hwAddr.String() != "")
}

func TestRssHashConfGet(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	var c RssConf
	err := pid.RssHashConfGet(&c)
	assert(t, err == nil, err)

	c.Key = bytes.Repeat([]byte{0x6d, 0x5a}, 20)
	err = pid.RssHashUpdate(&c)
	assert(t, err == nil, err)
}

func TestDevInfo(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	var c DevInfo
	err := pid.InfoGet(&c)
	assert(t, err == nil, err)

	assert(t, c.DriverName() == "net_null", c.DriverName())
}

func TestPortName(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	s, err := pid.Name()
	assert(t, err == nil, err)
	assert(t, s == "net_null0")

	otherPid, err := GetPortByName(s)
	assert(t, err == nil, err)
	assert(t, pid == otherPid)

	_, err = GetPortByName("some_name")
	assert(t, err == syscall.ENODEV)
}

func TestRssHashUpdate(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	conf := &RssConf{
		Key: make([]byte, 40),
		Hf:  0,
	}

	err := pid.RssHashUpdate(conf)
	assert(t, err == nil)
}
