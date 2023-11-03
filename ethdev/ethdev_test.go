package ethdev

import (
	"bytes"
	"errors"
	"syscall"
	"testing"

	"github.com/yerden/go-dpdk/eal"
)

func TestMACAddr(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	var a MACAddr
	err := pid.MACAddrGet(&a)
	assert(t, err == nil, err)
}

func TestRssHashConfGet(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	var c RssConf
	err := pid.RssHashConfGet(&c)
	assert(t, err == nil, err)

	c.Key = bytes.Repeat([]byte{0x6d, 0x5a}, 20)
	err = pid.RssHashUpdate(&c)
	assert(t, err == nil || errors.Is(err, syscall.ENOTSUP))
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
