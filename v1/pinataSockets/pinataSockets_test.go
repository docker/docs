package pinataSockets

import (
	"fmt"
	"testing"
)

func TestSocketPaths(t *testing.T) {

	var err error
	var str string

	fmt.Println(GetDBSocketPath())
	fmt.Println(GetOsxfsVolumeSocketPath())
	fmt.Println(GetSlirpSocketPath())
	fmt.Println(GetPortSocketPath())
	fmt.Println(GetDockerSocketPath())

	fmt.Println("cid:", 1, "- port:", 1)
	fmt.Println(GetVsockSocketPath(1, 1))

	fmt.Println("cid:", 4294967295, "- port:", 65535)
	fmt.Println(GetVsockSocketPath(4294967295, 65535))

	fmt.Println("cid:", 2, "- port:", 1524)
	fmt.Println(GetVsockSocketPath(2, 1524))

	fmt.Println("cid:", 3, "- port:", 1525)
	fmt.Println(GetVsockSocketPath(3, 1525))

	fmt.Println("cid:", 3, "- port:", 2376)
	fmt.Println(GetVsockSocketPath(3, 2376))

	fmt.Println("connect:")
	fmt.Println(GetVsockConnectSocketPath())

	fmt.Println("guest, port:", 1)
	str, err = GetVsockAliasSocketPath("guest", 1)
	if err == nil {
		fmt.Println(str)
	}

	fmt.Println("guest, port:", 65535)
	str, err = GetVsockAliasSocketPath("guest", 65535)
	if err == nil {
		fmt.Println(str)
	}
}
