package shell

import (
	"bytes"
	"fmt"
	"hash/crc64"
	"io"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/koud-fi/pkg/blob"

	"golang.org/x/sync/singleflight"
)

var (
	ctrlMap  sync.Map
	crcTable = crc64.MakeTable(crc64.ISO)
)

func Run(cmd string, args ...interface{}) blob.Blob {
	return blob.Func(func() (io.ReadCloser, error) {
		var (
			ctrl    = cmdCtrl(cmd)
			keyCrc  = crc64.New(crcTable)
			argStrs = make([]string, 0, len(args))
			stdin   blob.Blob
		)
		keyCrc.Write(*(*[]byte)(unsafe.Pointer(&cmd)))
		for _, arg := range args {
			switch v := arg.(type) {
			case blob.Blob:
				stdin = v
				continue
			}
			argStr := fmt.Sprint(arg)
			keyCrc.Write(*(*[]byte)(unsafe.Pointer(&argStr)))
			argStrs = append(argStrs, argStr)
		}
		key := strconv.FormatUint(keyCrc.Sum64(), 36)
		out, err, _ := ctrl.group.Do(key, func() (interface{}, error) {
			ctrl.throttle <- struct{}{}
			defer func() { <-ctrl.throttle }()

			cmd := exec.Command(cmd, argStrs...)
			if stdin != nil {
				rc, err := stdin.Open()
				if err != nil {
					return nil, err
				}
				defer rc.Close()
				cmd.Stdin = rc
			}
			out, err := cmd.CombinedOutput() // TODO: read output separately
			if err != nil {
				return []byte{}, fmt.Errorf("%w: %s", err,
					strings.TrimSpace(strings.ReplaceAll(string(out), "\r\n", " ")))
			}
			return out, nil
		})

		// TODO: avoid copying the full output to memory

		return io.NopCloser(bytes.NewReader(out.([]byte))), err
	})
}

type ctrl struct {
	group    singleflight.Group
	throttle chan struct{}
}

func cmdCtrl(cmd string) *ctrl {
	ctrlif, ok := ctrlMap.Load(cmd)
	if !ok {
		ctrlif, _ = ctrlMap.LoadOrStore(cmd, &ctrl{
			throttle: make(chan struct{}, runtime.NumCPU()),
		})
	}
	return ctrlif.(*ctrl)
}
