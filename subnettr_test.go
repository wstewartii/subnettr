package main

import "testing"

func TestSubnetter(t *testing.T) {
	subnet := subnettr("192", "255", 1)
	bcast := subnettr("1", "0", 2)
	if bcast != "255" {
		t.Errorf("Invalid broadcast return value for subnettr. Got %s, expected %s.", bcast, "255")
	}
	if subnet != "192" {
		t.Errorf("Invalid host return value for subnettr. Got %s, expected %s.", subnet, "192")
	}

}

func TestCidrToMask(t *testing.T) {
	nMask := cidrToMask("23")
	if nMask != "255.255.254.0" {
		t.Errorf("Invalid netmask return value for subnettr. Got %s, expected %s.", nMask, "255.255.254.0")
	}

}

func TestgetNetworkObject(t *testing.T) {
	netObj, err := getNetworkObject("192.168.1.5", "27")
	if err != nil {
		t.Errorf(err.Error())
	}
	if netObj.Subnet != "192.168.1.0" {
		t.Errorf("Invalid value for network address. Got %s, expected %s.", netObj.Subnet, "192.168.1.0")
	}
	if netObj.LastHostAddress != "192.168.1.30" {
		t.Errorf("Invalid address for last host on network. Got %s, expected %s.", netObj.LastHostAddress, "192.168.1.30")
	}
	if netObj.FirstHostAddress != "192.168.1.1" {
		t.Errorf("Invalid address for first host on network. Got %s, expected %s.", netObj.FirstHostAddress, "192.168.1.1")
	}
	if netObj.BroadcastAddress != "192.168.1.31" {
		t.Errorf("Invalid broadcast address. Got %s, expected %s.", netObj.BroadcastAddress, "192.168.1.31")
	}
	if netObj.SubnetMask != "255.255.255.224" {
		t.Errorf("Invalid netmask address. Got %s, expected %s.", netObj.SubnetMask, "255.255.255.224")
	}
}
