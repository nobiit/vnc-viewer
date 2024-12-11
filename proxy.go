package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func proxyVnc(c context.Context, addr string) (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return addr, fmt.Errorf("net.Listen: %v", err)
	}
	go func() {
		<-c.Done()
		if err = l.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Listener.Close: %v\n", err)
		}
	}()
	go func() {
		var (
			lConn net.Conn
			rConn net.Conn
		)
		for {
			lConn, err = l.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Listener.Accept: %v\n", err)
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				time.Sleep(time.Second)
				continue
			}

			for {
				rConn, err = net.Dial("tcp", addr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "net.Dial: %v\n", err)
					if strings.Contains(err.Error(), "connect: connection refused") {
						time.Sleep(time.Second)
						continue
					}
					_ = l.Close()
				}
				break
			}

			cErr := make(chan error, 1)
			go func() {
				_, err = io.Copy(lConn, rConn)
				cErr <- err
			}()
			go func() {
				_, err = io.Copy(rConn, lConn)
				cErr <- err
			}()
			go func() {
				<-c.Done()
				cErr <- c.Err()
			}()
			if err = errors.Join(<-cErr, lConn.Close(), rConn.Close()); err != nil {
				fmt.Fprintf(os.Stderr, "Connection: %v\n", err)
				continue
			}
		}
	}()
	return l.Addr().String(), nil
}
