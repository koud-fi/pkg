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
	DefaultTLSAddr      = ":https"

	autoCertCacheDir = ".autocert-cache"
)

type listenConfig struct {
	addr      string
	tlsConfig *tls.Config
}

type ListenOption func(*listenConfig)

func Addr(addr string) ListenOption {
	return func(c *listenConfig) { c.addr = addr }
}

func TLS(c *tls.Config) ListenOption {
	return func(lc *listenConfig) {
		if c == nil {
			return
		}
		if lc.addr == DefaultAddr {
			lc.addr = DefaultTLSAddr
		}
		lc.tlsConfig = c
	}
}

func Listen(h http.Handler, opt ...ListenOption) {
	c := listenConfig{
		addr: DefaultAddr,
	}
	for _, opt := range opt {
		opt(&c)
	}
	runServer(c.addr, h, c.tlsConfig)
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
