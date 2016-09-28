// +build darwin
package client

import (
	"errors"
)

func (c *EngineClient) CheckKernelVersion() error {
	return errors.New("CheckKernelVersion not implemented on Darwin")
}
