package internal

import (
	"net/http/httputil"
	"net/url"
)

type Server struct {
	_url  *url.URL
	proxy *httputil.ReverseProxy
}

func CreateServer(u *url.URL) *Server {
	return &Server{
		_url:  u,
		proxy: httputil.NewSingleHostReverseProxy(u),
	}
}
