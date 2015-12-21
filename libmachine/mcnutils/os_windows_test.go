package mcnutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVerOutput(t *testing.T) {
	output := `

Microsoft Windows [version 6.3.9600]

`

	assert.Equal(t, "Microsoft Windows [version 6.3.9600]", parseVerOutput(output))
}

func TestParseSystemInfoOutput(t *testing.T) {
	output := `
Host Name:                 DESKTOP-3A5PULA
OS Name:                   Microsoft Windows 10 Enterprise
OS Version:                10.0.10240 N/A Build 10240
OS Manufacturer:           Microsoft Corporation
OS Configuration:          Standalone Workstation
OS Build Type:             Multiprocessor Free
Registered Owner:          Windows User
`

	assert.Equal(t, "10.0.10240 N/A Build 10240", parseOutput(output))
}
