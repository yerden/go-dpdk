package eal

import (
	"fmt"
	"strings"
)

// single EAL command line option
type cmdOption struct {
	option string
	value  *string
}

// create either option-value or single option
func newOption(tokens ...string) cmdOption {
	switch len(tokens) {
	case 1:
		return cmdOption{option: tokens[0]}
	case 2:
		fallthrough
	default:
		b := tokens[1]
		return cmdOption{option: tokens[0], value: &b}
	}
}

// options array for EAL
type ealOptions struct {
	opts []cmdOption
}

func (p *ealOptions) find(option string) *cmdOption {
	for i := range p.opts {
		if opt := &p.opts[i]; opt.option == option {
			return opt
		}
	}
	return nil
}

// adds new options or replaces the option to EAL if it was already
// given
func (p *ealOptions) replace(in cmdOption) {
	if opt := p.find(in.option); opt != nil {
		*opt = in
	} else {
		p.opts = append(p.opts, in)
	}
}

// add the option to EAL
func (p *ealOptions) add(in cmdOption) {
	p.opts = append(p.opts, in)
}

// convert options array into argv array for rte_eal_init.
func (p *ealOptions) argv() []string {
	s := []string{}
	for _, opt := range p.opts {
		if opt.option == "" {
			panic("incorrect options code")
		} else if opt.value == nil {
			s = append(s, opt.option)
		} else if opt.option[:2] == "--" {
			s = append(s, opt.option+"="+*opt.value)
		} else {
			s = append(s, opt.option)
			s = append(s, *opt.value)
		}
	}
	return s
}

func (p *ealOptions) String() string {
	return fmt.Sprintln(p.argv())
}

// Option represents an option for EAL initialization. Options may be
// combined in array. Some options may be specified several times in a
// row, some may only happen once.
type Option struct {
	f func(*ealOptions)
}

// OptLcores sets mask of the cores to run on. Please note that this
// is a mandatory option.
func OptLcores(mask Set) Option {
	opt := newOption("-c", SetToHex(mask, MaxLcore))
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptMasterLcore specifies core ID that is used as master.
func OptMasterLcore(n int) Option {
	opt := newOption("--master-lcore", fmt.Sprintf("%d", n))
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptServiceLcores specifies mask of cores to be used as service
// cores.
func OptServiceLcores(mask Set) Option {
	opt := newOption("-s", SetToHex(mask, MaxLcore))
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptBlacklistDev blacklists a PCI device to prevent EAL from using
// it. Multiple option instances are allowed. This option negates
// OptWhitelistDev.
func OptBlacklistDev(dev string) Option {
	opt := newOption("--pci-blacklist", dev)
	return Option{func(p *ealOptions) { p.add(opt) }}
}

// OptWhitelistDev adds a PCI device in white list. Multiple option
// instances are allowed. This option negates OptBlacklistDev.
func OptWhitelistDev(dev string) Option {
	opt := newOption("--pci-whitelist", dev)
	return Option{func(p *ealOptions) { p.add(opt) }}
}

// OptFilePrefix specifies a different shared data file prefix for a
// DPDK process. This option allows running multiple independent DPDK
// primary/secondary processes under different prefixes.
func OptFilePrefix(prefix string) Option {
	opt := newOption("--file-prefix", prefix)
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptBaseVirtAddr attempts to use a different starting address for all
// memory maps of the primary DPDK process. This can be helpful if
// secondary processes cannot start due to conflicts in address map.
func OptBaseVirtAddr(addr uintptr) Option {
	opt := newOption("--base-virtaddr", fmt.Sprintf("0x%x", addr))
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptLoadExternalPath loads external drivers. An argument can be a
// single shared object file, or a directory containing multiple
// driver shared objects. Multiple option instances are allowed.
func OptLoadExternalPath(path string) Option {
	opt := newOption("-d", path)
	return Option{func(p *ealOptions) { p.add(opt) }}
}

// OptNoPCI disables PCI bus.
func OptNoPCI() Option {
	opt := newOption("--no-pci")
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptNoHuge uses anonymous memory instead of hugepages (implies no
// secondary process support).
func OptNoHuge() Option {
	opt := newOption("--no-huge")
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptProcType sets the type of the current process. It can either be
// ProcPrimary, ProcSecondary or ProcAuto.
func OptProcType(typ int) Option {
	var opt cmdOption
	switch typ {
	default:
		fallthrough
	case ProcPrimary:
		opt = newOption("--proc-type", "primary")
	case ProcSecondary:
		opt = newOption("--proc-type", "secondary")
	case ProcAuto:
		opt = newOption("--proc-type", "auto")
	}
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptMemoryChannels sets the number of memory channels to use.
func OptMemoryChannels(n int) Option {
	opt := newOption("-n", fmt.Sprintf("%d", n))
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptInMemory instructs not to create any shared data structures and run
// entirely in memory.  Implies --no-shconf and (if applicable) --huge-unlink.
func OptInMemory() Option {
	opt := newOption("--in-memory")
	return Option{func(p *ealOptions) { p.replace(opt) }}
}

// OptSocketMemory preallocates specified amounts of memory per
// socket. The number is specified as megabytes with each value for
// a socket.
func OptSocketMemory(mem ...int) Option {
	var s []string
	for _, m := range mem {
		s = append(s, fmt.Sprintf("%d", m))
	}
	opt := newOption("--socket-mem", strings.Join(s, ","))
	return Option{func(p *ealOptions) { p.replace(opt) }}
}
