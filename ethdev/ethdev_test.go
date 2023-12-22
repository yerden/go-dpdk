package ethdev

import (
	"bytes"
	"errors"
	"math"
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
	assert(t, c.DevFlags() == 64, c.DevFlags())
	assert(t, c.InterfaceName() == "", c.InterfaceName())
	assert(t, c.MaxMTU() == math.MaxUint16, c.MaxMTU())
	assert(t, c.MinMTU() == 46, c.MinMTU())
	assert(t, c.MaxRxPktLen() == math.MaxUint32, c.MaxRxPktLen())
	assert(t, c.MaxRxQueues() == 1024, c.MaxRxQueues())
	assert(t, c.MaxTxQueues() == 1024, c.MaxTxQueues())
	assert(t, c.MinRxBufSize() == 0, c.MinRxBufSize())
	assert(t, c.NbRxQueues() == 1, c.NbRxQueues())
	assert(t, c.NbTxQueues() == 1, c.NbTxQueues())
	assert(t, c.RetaSize() == 128, c.RetaSize())
}

func TestMTU(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	mtu, err := pid.GetMTU()
	assert(t, mtu == 1500, mtu)
	assert(t, err == nil, err)

	err = pid.DevConfigure(1, 1)
	assert(t, err == nil, err)

	err = pid.SetMTU(math.MaxUint16)
	assert(t, err == nil, err)
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

func TestOptRxMode(t *testing.T) {
	opt := OptRxMode(RxMode{
		MqMode:       1,
		MTU:          2,
		SplitHdrSize: 3,
		Offloads:     4,
	})
	cfg := &ethConf{}
	opt.f(cfg)

	assert(t, cfg.conf.rxmode.mq_mode == 1)
	assert(t, cfg.conf.rxmode.mtu == 2)
	assert(t, cfg.conf.rxmode.offloads == 4)
}

func TestOptIntrConf(t *testing.T) {
	opt := OptIntrConf(IntrConf{})

	cfg := &ethConf{}
	opt.f(cfg)

	flags := readIntrConf(&cfg.conf.intr_conf)
	assert(t, flags[0] == 0)
	assert(t, flags[1] == 0)
	assert(t, flags[2] == 0)

	opt = OptIntrConf(IntrConf{
		LSC: true,
	})

	cfg = &ethConf{}
	opt.f(cfg)

	flags = readIntrConf(&cfg.conf.intr_conf)
	assert(t, flags[0] == 1)
	assert(t, flags[1] == 0)
	assert(t, flags[2] == 0)

	opt = OptIntrConf(IntrConf{
		LSC: true,
		RXQ: true,
	})

	cfg = &ethConf{}
	opt.f(cfg)

	flags = readIntrConf(&cfg.conf.intr_conf)
	assert(t, flags[0] == 1)
	assert(t, flags[1] == 1)
	assert(t, flags[2] == 0)

	opt = OptIntrConf(IntrConf{
		RXQ: true,
		RMV: true,
	})

	cfg = &ethConf{}
	opt.f(cfg)

	flags = readIntrConf(&cfg.conf.intr_conf)
	assert(t, flags[0] == 0)
	assert(t, flags[1] == 1)
	assert(t, flags[2] == 1)
}

func TestLSC(t *testing.T) {
	eal.InitOnceSafe("test", 4)

	pid := Port(0)

	// net_null0, no support for LSC
	var info DevInfo
	err := pid.InfoGet(&info)
	assert(t, err == nil, err)
	assert(t, info.DevFlags().IsIntrLSC() == false)

	err = pid.RegisterCallbackLSC()
	assert(t, err == nil, err)

	err = pid.UnregisterCallbackLSC()
	assert(t, err == nil, err)

	RegisterTelemetryLSC("/ethdev/lsc")
}
