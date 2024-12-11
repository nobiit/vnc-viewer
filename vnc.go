package main

import (
	"runtime"
)

var Path string

func init() {
	switch runtime.GOOS {
	case "darwin":
		Path = "/Applications/VNC Viewer.app/Contents/MacOS/vncviewer"
	}
}
