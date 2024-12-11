package main

import (
	"context"
	"errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"sync"
)

func runAction(c *cli.Context) error {
	var errs []error
	w := new(sync.WaitGroup)
	for _, item := range c.Args().Slice() {
		ctx, cancel := context.WithCancel(c.Context)
		addr, err := proxyVnc(ctx, item)
		if err != nil {
			cancel()
			errs = append(errs, err)
			continue
		}
		s := exec.Command(Path, "-WarnUnencrypted=false", addr)
		s.Stdout = os.Stdout
		s.Stderr = os.Stderr
		w.Add(1)
		go func() {
			defer w.Done()
			defer cancel()
			if err = s.Run(); err != nil {
				errs = append(errs, err)
			}
		}()
	}
	w.Wait()
	return errors.Join(errs...)
}
