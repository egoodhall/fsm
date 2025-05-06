package fsm

import (
	"errors"
	"flag"
)

type GeneratorOptions struct {
	Out string
	Pkg string
}

func ParseFlags() (input string, opts GeneratorOptions, err error) {
	flag.StringVar(&opts.Out, "out", "", "output directory")
	flag.StringVar(&opts.Pkg, "pkg", "", "package name")
	flag.Parse()

	if opts.Out == "" {
		return "", opts, errors.New("output directory is required")
	} else if opts.Pkg == "" {
		return "", opts, errors.New("package name is required")
	}

	return flag.Arg(0), opts, nil
}
