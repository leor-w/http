package http

import "github.com/leor-w/kid/server"

func newOptions(opt ...server.Option) server.Options {
	opts := server.Options{}
	for _, o := range opt {
		o(&opts)
	}
	if len(opts.Address) == 0 {
		opts.Address = ":8080"
	}
	return opts
}
