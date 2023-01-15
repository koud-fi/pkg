package serve

import (
	"crypto/tls"
	"os"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"
)

func AutoCert(email string, host string) *tls.Config {
	if email == "" || host == "" {
		return nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("serve.AutoCert: unable to resolve home dir")
	}
	return (&autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(host),

		Email: email,
		Cache: autocert.DirCache(filepath.Join(homeDir, autoCertCacheDir)),
	}).TLSConfig()
}
