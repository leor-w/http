package http

// httpHandler http.Handler 的包装, 只有一个接口, 返回 handler
type httpHandler struct {
	handler interface{}
}

func (h *httpHandler) Handler() interface{} {
	return h.handler
}
