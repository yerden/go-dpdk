package bpf

import "log"

func ExampleLoad() {
	prm := &Prm{
		Insns: NewInsns([]Insn{
			// ...
		}),
		XSyms: NewXSyms([]XSym{
			// ...
		}),
		ProgArg: &ArgPtr{
			Size: 1024,
		},
	}

	bpf, err := Load(prm)
	if err != nil {
		log.Fatalln("unable to load program:", err)
	}

	_ = bpf
}
