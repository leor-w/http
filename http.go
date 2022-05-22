package http

import (
	"crypto/tls"
	"errors"
	"github.com/leor-w/kid/server"
	"net"
	"net/http"
	"sync"
)

type httpServer struct {
	sync.Mutex
	opts    server.Options
	handler server.Handler
	exit    chan error
}

// Init Server 初始化 Option
func (h *httpServer) Init(options ...server.Option) {
	h.Lock()
	for _, o := range options {
		o(&h.opts)
	}
	h.Unlock()
}

// Options 设置 Options
func (h *httpServer) Options() server.Options {
	h.Lock()
	opts := h.opts
	h.Unlock()
	return opts
}

// Handle 设置 Handler 必须是 go net 包下的 http.Handler 接口的实现, 否则将返回错误
func (h *httpServer) Handle(handler server.Handler) error {
	if _, ok := handler.Handler().(http.Handler); !ok {
		return errors.New("handler 必须为 http.Handler")
	}
	h.Lock()
	h.handler = handler
	h.Unlock()
	return nil
}

// NewHandler 创建一个 http.Handler 的包装
func (h *httpServer) NewHandler(handler interface{}) server.Handler {
	return &httpHandler{
		handler: handler,
	}
}

// Start 启动 http 服务
func (h *httpServer) Start() error {
	h.Lock()
	opts := h.opts
	handler := h.handler
	h.Unlock()
	var (
		ln  net.Listener
		err error
	)

	if opts.TLSConfig != nil {
		ln, err = tls.Listen("tcp", opts.Address, opts.TLSConfig)
	} else {
		ln, err = net.Listen("tcp", opts.Address)
	}
	if err != nil {
		return err
	}
	// TODO: 加入启动提示信息

	h.Lock()
	h.opts.Address = ln.Addr().String()
	h.Unlock()

	hd, ok := handler.Handler().(http.Handler)
	if !ok {
		return errors.New("server 必须为 http.Handler")
	}

	go http.Serve(ln, hd)

	go func() {
	Loop:
		for {
			select {
			// 等待发送退出信号
			case <-h.exit:
				break Loop
			}
		}
		_ = ln.Close()
	}()
	return nil
}

func (h *httpServer) Stop() error {
	err := errors.New("")
	h.exit <- err
	return err
}

func newServer(opts ...server.Option) *httpServer {
	return &httpServer{
		opts: newOptions(opts...),
		exit: make(chan error),
	}
}

func NewServer(opts ...server.Option) *httpServer {
	return newServer(opts...)
}
