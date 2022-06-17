package serve

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

const (
	DefaultReadTimeout  = time.Second * 30
	DefaultWriteTimeout = time.Second * 60
	DefaultIdleTimeout  = time.Second * 120
	DefaultAddr         = ":http"
	//DefaultTLSAddr = ":https"
)

type listenConfig struct {
	addr string
}

type ListenOption func(*listenConfig)

func Addr(addr string) ListenOption {
	return func(c *listenConfig) { c.addr = addr }
}

func Listen(h http.Handler, opt ...ListenOption) {
	c := listenConfig{
		addr: DefaultAddr,
	}
	for _, opt := range opt {
		opt(&c)
	}

	// TODO

	runServer(c.addr, h, nil)
}

func runServer(addr string, h http.Handler, tlsConf *tls.Config) {
	if h == nil {
		h = http.DefaultServeMux
	}
	s := http.Server{
		Addr:      addr,
		Handler:   h,
		TLSConfig: tlsConf,

		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,

		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	log.Println("LISTEN", addr)
	log.Fatal(s.ListenAndServe())
}
