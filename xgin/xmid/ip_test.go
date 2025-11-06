package xmid

import (
	"testing"
)

func Test_ip2int(t *testing.T) {
	ip := "192.168.32.21"

	iv, err := ip2int(ip)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(iv)
	}
}

func Test_parseRange(t *testing.T) {
	ip := "192.168.32.21/24"

	iv, err := parseIpRange(ip)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(iv)
	}
}

func Test_CheckIp(t *testing.T) {
	ip := "192.168.31.21"
	rg := "192.168.32.0/24"

	iv, _ := ip2int(ip)
	ir, _ := parseIpRange(rg)

	t.Log(iv&ir == ir)
}
