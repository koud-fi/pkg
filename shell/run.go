package shell

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc64"
	"io"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"unsafe"

	"github.com/koud-fi/pkg/blob"

	"golang.org/x/sync/singleflight"
)

var (
	ctrlMap  sync.Map
	crcTable = crc64.MakeTable(crc64.ISO)
)

func Run(ctx context.Context, cmd string, args ...interface{}) blob.Blob {
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
			var (
				outBuf = bytes.NewBuffer(nil)
				errBuf = bytes.NewBuffer(nil)
			)
			cmd.Stdout = outBuf
			cmd.Stderr = errBuf

			if err := cmd.Run(); err != nil {
				if ctxErr := ctx.Err(); ctxErr != nil {
					return nil, ctxErr
				}
				switch err := err.(type) {
				case *exec.ExitError:
					msg, _ := errBuf.ReadString('\n')

					// TODO: support multi-line error output

					return nil, fmt.Errorf("exit status %d: %s", err.ExitCode(), msg)
				}
				return nil, err
			}
			return outBuf.Bytes(), nil
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
