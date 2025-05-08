package serve

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/koud-fi/pkg/netx/httpmiddleware/logger"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
)

type HostPolicy = autocert.HostPolicy

func TLS(h http.Handler, hostPolicy HostPolicy) error {
	const addr = "0.0.0.0:443"

	// TODO: better configuration
	// TODO: don't automatically wrap with loggers

	var (
		certMgr = autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(".certs"), // TODO: make configurable
			HostPolicy: hostPolicy,
		}
		g errgroup.Group
	)
	g.Go(func() error {
		return http.ListenAndServe("0.0.0.0:80", logger.Wrap(certMgr.HTTPHandler(nil)))
	})
	g.Go(func() error {
		srv := http.Server{
			Addr:      addr,
			Handler:   h,
			TLSConfig: newTLSConfig(&certMgr),
		}
		return srv.ListenAndServeTLS("", "")
	})
	log.Println("listening on", addr)

	return g.Wait()
}

func HostWhitelist(hosts ...string) HostPolicy {
	return autocert.HostWhitelist(hosts...)
}

func HostSuffix(suffix string) HostPolicy {
	return func(_ context.Context, host string) error {
		if !strings.HasSuffix(host, suffix) {
			return fmt.Errorf("serve.TLS: host %q doesn't match suffix %q", host, suffix)
		}
		return nil
	}
}

func newTLSConfig(certMgr *autocert.Manager) *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		GetCertificate: certMgr.GetCertificate,
		NextProtos: []string{
			"h2", "http/1.1",
			acme.ALPNProto,
		},
	}
}
