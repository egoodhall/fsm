package fsm

import "flag"

type Options struct {
	Out string
}

func ParseOptions() (input string, opts Options, err error) {
	flag.StringVar(&opts.Out, "out", "", "output directory")
	flag.Parse()
	return flag.Arg(0), opts, nil
}
