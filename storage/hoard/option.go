package hoard

import (
	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file"
	"github.com/koud-fi/pkg/pk"
)

type config struct {
	src          blob.Mux
	idSalt       pk.Salt
	fileAttrOpts []file.Option
}

type Option func(*config)

func Source(src blob.Mux) Option             { return func(c *config) { c.src = src } }
func IDSalt(salt pk.Salt) Option             { return func(c *config) { c.idSalt = salt } }
func FileAttrOpts(opt ...file.Option) Option { return func(c *config) { c.fileAttrOpts = opt } }
