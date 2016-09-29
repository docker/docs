// +build windows

package main

import (
	"net"

	"github.com/Microsoft/go-winio"
)

func Dial(path string) (net.Conn, error) {
	return winio.DialPipe(path, nil)
}

func GetDefaultDBPath() string {
	return `\\.\pipe\dockerDataBase`
}
