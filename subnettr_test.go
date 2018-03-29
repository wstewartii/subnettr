package main

import "testing"

func TestCidrToMask(t *testing.T) {
	nMask := cidrToMask("23")
	if nMask != "255.255.254.0" {
		t.Errorf("Invalid netmask return value for subnettr. Got %s, expected %s.", nMask, "255.255.254.0")
	}

}

func TestGetNetworkObject(t *testing.T) {
	netObj, err := getNetworkObject("192.168.1.5", "27")
	if err != nil {
		t.Errorf(err.Error())
	}
	if netObj.NetworkID.String() != "192.168.1.0" {
		t.Errorf("Invalid value for network address. Got %s, expected %s.", netObj.NetworkID, "192.168.1.0")
	}
	if netObj.BroadcastAddress.String() != "192.168.1.31" {
		t.Errorf("Invalid broadcast address. Got %s, expected %s.", netObj.BroadcastAddress, "192.168.1.31")
	}
	if netObj.SubnetMask.String() != "255.255.255.224" {
		t.Errorf("Invalid netmask address. Got %s, expected %s.", netObj.SubnetMask, "255.255.255.224")
	}
	if netObj.UsableHostAddresses != 30 {
		t.Errorf("Invalid value for usable host addresses. Got %f, expected %f.", netObj.UsableHostAddresses, float64(30))
	}
}
