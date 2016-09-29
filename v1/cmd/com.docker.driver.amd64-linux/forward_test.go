package main

import (
	"strings"
	"testing"
)

func TestLeaseParser(t *testing.T) {
	exampleLeases := `
	{
        ip_address=192.168.64.130
        hw_address=1,6:6d:15:4f:cf:32
        identifier=1,6:6d:15:4f:cf:32
        lease=0x56d86f92
  }
	{
	        ip_address=192.168.64.129
	        hw_address=1,fe:82:1b:45:cf:32
	        identifier=1,fe:82:1b:45:cf:32
	        lease=0x56d852b9
	}
	{
	        name=docker
	        ip_address=192.168.64.128
	        hw_address=1,da:c6:98:71:cf:32
	        identifier=1,da:c6:98:71:cf:32
	        lease=0x56d730a8
	}
	{
	        name=docker
	        ip_address=192.168.64.127
	        hw_address=1,76:1b:f:f:cf:32
	        identifier=1,76:1b:f:f:cf:32
	        lease=0x56d72e2f
	}
	{
	        ip_address=192.168.64.126
	        hw_address=1,26:fd:f0:3:cf:32
	        identifier=1,26:fd:f0:3:cf:32
	        lease=0x56d72732
	}
	{
	        name=docker
	        ip_address=192.168.64.125
	        hw_address=1,4e:13:89:15:cf:32
	        identifier=1,4e:13:89:15:cf:32
	        lease=0x56d6f1d8
	}
	{
	        name=docker
	        ip_address=192.168.64.124
	        hw_address=1,2:d5:28:5d:cf:32
	        identifier=1,2:d5:28:5d:cf:32
	        lease=0x56d588a9
	}
	{
	        name=docker
	        ip_address=192.168.64.123
	        hw_address=1,ca:a7:7:65:cf:32
	        identifier=1,ca:a7:7:65:cf:32
	        lease=0x56d5871d
	}
	{
	        name=docker
	        ip_address=192.168.64.122
	        hw_address=1,ae:2f:b:61:cf:32
	        identifier=1,ae:2f:b:61:cf:32
	        lease=0x56d5793e
	}
	{
	        name=docker
	        ip_address=192.168.64.121
	        hw_address=1,ca:f8:5b:2:cf:32
	        identifier=1,ca:f8:5b:2:cf:32
	        lease=0x56d57863
	}
	{
	        name=docker
	        ip_address=192.168.64.120
	        hw_address=1,2e:89:ac:46:cf:32
	        identifier=1,2e:89:ac:46:cf:32
	        lease=0x56d1ed83
	}
	{
	        name=docker
	        ip_address=192.168.64.119
	        hw_address=1,c2:8e:82:35:cf:32
	        identifier=1,c2:8e:82:35:cf:32
	        lease=0x56d1ea94
	}
	{
	        name=docker
	        ip_address=192.168.64.118
	        hw_address=1,8a:98:f4:15:cf:32
	        identifier=1,8a:98:f4:15:cf:32
	        lease=0x56d1ea83
	}
	{
	        name=docker
	        ip_address=192.168.64.117
	        hw_address=1,82:5d:3b:32:cf:32
	        identifier=1,82:5d:3b:32:cf:32
	        lease=0x56d1e9e7
	}
	{
	        name=docker
	        ip_address=192.168.64.116
	        hw_address=1,ae:3f:94:34:cf:32
	        identifier=1,ae:3f:94:34:cf:32
	        lease=0x56d1e9d0
	}
	{
	        name=docker
	        ip_address=192.168.64.115
	        hw_address=1,16:73:fc:62:cf:32
	        identifier=1,16:73:fc:62:cf:32
	        lease=0x56d1e19d
	}
	{
	        name=docker
	        ip_address=192.168.64.114
	        hw_address=1,1a:6c:c0:b:cf:32
	        identifier=1,1a:6c:c0:b:cf:32
	        lease=0x56d1e18b
	}
	{
	        name=docker
	        ip_address=192.168.64.113
	        hw_address=1,7e:ad:6a:6:cf:32
	        identifier=1,7e:ad:6a:6:cf:32
	        lease=0x56d1e0d9
	}
	{
	        name=docker
	        ip_address=192.168.64.112
	        hw_address=1,6e:79:14:67:cf:32
	        identifier=1,6e:79:14:67:cf:32
	        lease=0x56d1e0c7
	}
	{
	        name=docker
	        ip_address=192.168.64.111
	        hw_address=1,fe:84:ed:17:cf:32
	        identifier=1,fe:84:ed:17:cf:32
	        lease=0x56d1ddc0
	}
	{
	        name=docker
	        ip_address=192.168.64.110
	        hw_address=1,36:ca:7a:2f:cf:32
	        identifier=1,36:ca:7a:2f:cf:32
	        lease=0x56d1dc2c
	}
	{
	        name=docker
	        ip_address=192.168.64.109
	        hw_address=1,d6:da:40:2b:cf:32
	        identifier=1,d6:da:40:2b:cf:32
	        lease=0x56d1db7c
	}
	{
	        name=docker
	        ip_address=192.168.64.108
	        hw_address=1,72:d3:8d:77:cf:32
	        identifier=1,72:d3:8d:77:cf:32
	        lease=0x56d1d9f9
	}
	{
	        name=docker
	        ip_address=192.168.64.107
	        hw_address=1,8a:8c:50:33:cf:32
	        identifier=1,8a:8c:50:33:cf:32
	        lease=0x56d1d8a3
	}
	{
	        name=docker
	        ip_address=192.168.64.106
	        hw_address=1,c6:8:79:18:cf:32
	        identifier=1,c6:8:79:18:cf:32
	        lease=0x56d1d79e
	}
	{
	        name=docker
	        ip_address=192.168.64.105
	        hw_address=1,6e:5b:73:7:cf:32
	        identifier=1,6e:5b:73:7:cf:32
	        lease=0x56d1d73f
	}
	{
	        name=docker
	        ip_address=192.168.64.104
	        hw_address=1,ae:4a:f3:29:cf:32
	        identifier=1,ae:4a:f3:29:cf:32
	        lease=0x56d1d678
	}
	{
	        name=docker
	        ip_address=192.168.64.4
	        hw_address=1,ba:a6:b0:26:cf:32
	        identifier=1,ba:a6:b0:26:cf:32
	        lease=0x566c3568
	}
	{
	        name=docker
	        ip_address=192.168.64.3
	        hw_address=1,3a:78:7e:49:cf:32
	        identifier=1,3a:78:7e:49:cf:32
	        lease=0x568e41e9
	}
	{
	        name=docker
	        ip_address=192.168.64.2
	        hw_address=1,ba:30:6a:62:cf:32
	        identifier=1,ba:30:6a:62:cf:32
	        lease=0x566ab876
	}
	`

	var tests = []struct {
		mac string
		ip  string
	}{
		{"06:6D:15:4f:cf:32", "192.168.64.130"}, // first, note the leading 0 and capital D
		{"fe:82:1b:45:cf:32", "192.168.64.129"}, // second
		{"da:c6:98:71:cf:32", "192.168.64.128"}, // third
		{"3a:78:7e:49:cf:32", "192.168.64.3"},   // second last
		{"ba:30:6a:62:cf:32", "192.168.64.2"},   // last
	}
	for _, test := range tests {
		mac, err := parseMAC(test.mac)
		if err != nil {
			t.Errorf("Failed to parse MAC %s: %#v", test.mac, err)
			continue
		}
		ip, err := readDHCPLease(strings.NewReader(exampleLeases), mac)
		if err != nil {
			t.Errorf("Failed to find MAC %s: %#v", test.mac, err)
			continue
		}
		if ip != test.ip {
			t.Errorf("Got wrong IP, expected %s, got %s", test.ip, ip)
			continue
		}
	}
}
