package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

var Cmd = &cli.App{
	Name:   "vnc-viewer",
	Action: runAction,
}

func main() {
	if err := Cmd.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(Cmd.ErrWriter, err)
		cli.OsExiter(1)
	}
}
