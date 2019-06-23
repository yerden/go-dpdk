package eal

import (
	"fmt"
	"strings"
)

var (
	optAppend  = false
	optReplace = !optAppend
)

// single EAL command line option
type cmdOption struct {
	option string
	value  *string
}

// create either option-value or single option
func newOption(replace bool, tokens ...string) Option {
	var x cmdOption

	x.option = tokens[0]
	if len(tokens) > 1 {
		b := tokens[1]
		x.value = &b
	}

	return Option{func(p *ealOptions) { p.add(x, replace) }}
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
func (p *ealOptions) add(in cmdOption, replace bool) {
	if replace == optReplace {
		if opt := p.find(in.option); opt != nil {
			*opt = in
			return
		}
	}
	p.opts = append(p.opts, in)
}

func optFlag(key string, replace bool) Option {
	return newOption(replace, key)
}

func optInteger(key string, value interface{}, replace bool) Option {
	return newOption(replace, key, fmt.Sprintf("%d", value))
}

func optString(key string, value string, replace bool) Option {
	return newOption(replace, key, value)
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
	return optString("-c", SetToHex(mask, MaxLcore), optReplace)
}

// OptMasterLcore specifies core ID that is used as master.
func OptMasterLcore(n int) Option {
	return optInteger("--master-lcore", n, optReplace)
}

// OptServiceLcores specifies mask of cores to be used as service
// cores.
func OptServiceLcores(mask Set) Option {
	return optString("-s", SetToHex(mask, MaxLcore), optReplace)
}

// OptBlacklistDev blacklists a PCI device to prevent EAL from using
// it. Multiple option instances are allowed. This option negates
// OptWhitelistDev.
func OptBlacklistDev(dev string) Option {
	return optString("--pci-blacklist", dev, optAppend)
}

// OptWhitelistDev adds a PCI device in white list. Multiple option
// instances are allowed. This option negates OptBlacklistDev.
func OptWhitelistDev(dev string) Option {
	return optString("--pci-whitelist", dev, optAppend)
}

// OptFilePrefix specifies a different shared data file prefix for a
// DPDK process. This option allows running multiple independent DPDK
// primary/secondary processes under different prefixes.
func OptFilePrefix(prefix string) Option {
	return optString("--file-prefix", prefix, optReplace)
}

// OptMemory specifies amount of memory to preallocate at startup.
func OptMemory(n int) Option {
	return optInteger("-m", n, optReplace)
}

// OptBaseVirtAddr attempts to use a different starting address for all
// memory maps of the primary DPDK process. This can be helpful if
// secondary processes cannot start due to conflicts in address map.
func OptBaseVirtAddr(addr uintptr) Option {
	return optInteger("--base-virtaddr", addr, optReplace)
}

// OptLoadExternalPath loads external drivers. An argument can be a
// single shared object file, or a directory containing multiple
// driver shared objects. Multiple option instances are allowed.
func OptLoadExternalPath(path string) Option {
	return optString("-d", path, optAppend)
}

// OptProcType sets the type of the current process. It can either be
// ProcPrimary, ProcSecondary or ProcAuto.
func OptProcType(typ int) Option {
	switch typ {
	default:
		fallthrough
	case ProcPrimary:
		return optString("--proc-type", "primary", optReplace)
	case ProcSecondary:
		return optString("--proc-type", "secondary", optReplace)
	case ProcAuto:
		return optString("--proc-type", "auto", optReplace)
	}
}

// OptMemoryChannels sets the number of memory channels to use.
func OptMemoryChannels(n int) Option {
	return optInteger("-n", n, optReplace)
}

// OptSocketMemory preallocates specified amounts of memory per
// socket. The number is specified as megabytes with each value for
// a socket.
func OptSocketMemory(mem ...int) Option {
	var s []string
	for _, m := range mem {
		s = append(s, fmt.Sprintf("%d", m))
	}
	return optString("--socket-mem", strings.Join(s, ","), optReplace)
}

var (
	// OptNoPCI disables PCI bus.
	OptNoPCI = optFlag("--no-pci", optReplace)

	// OptNoHuge uses anonymous memory instead of hugepages (implies
	// no secondary process support).
	OptNoHuge = optFlag("--no-huge", optReplace)

	// OptInMemory instructs not to create any shared data structures
	// and run entirely in memory.  Implies --no-shconf and (if
	// applicable) --huge-unlink.
	OptInMemory = optFlag("--in-memory", optReplace)
)

// OptArgs parses array of Options and return argv array for
// rte_eal_init() call.
func OptArgs(opts []Option) []string {
	p := ealOptions{}
	for _, opt := range opts {
		opt.f(&p)
	}

	return p.argv()
}
