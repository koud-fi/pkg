package cli

import (
	"context"
	"os"
	"strings"

	"github.com/koud-fi/pkg/proc"
	"github.com/koud-fi/pkg/proc/router"
)

type CLI struct {
	r router.Router
}

func New(r router.Router) CLI {
	return CLI{r: r}
}

func (c CLI) Run(ctx context.Context, args ...string) error {
	var cmd string
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}
	return writeOutput(c.r.Invoke(ctx, cmd, Params(args...)))
}

func Run(ctx context.Context, p proc.Proc, args ...string) error {
	return writeOutput(p.Invoke(ctx, Params(args...)))
}

func writeOutput(out any, err error) error {
	if err != nil {
		return err
	}
	return proc.WriteOutput(os.Stdout, out)
}

func Params(args ...string) proc.Params {
	var (
		m       = make(proc.ParamMap)
		currKey string
	)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			currKey = strings.TrimLeft(arg, "-")
			if _, ok := m[currKey]; !ok {
				m[currKey] = make([]string, 0)
			}
		} else if currKey != "" {
			m[currKey] = append(m[currKey], arg)
		}
	}
	return m
}
