package main

import "strconv"
import "testing"

func TestSubnetter(t *testing.T) {
	subnet := strconv.FormatUint(uint64(subnettr("192", "255", 1)),10)
	bcast := strconv.FormatUint(uint64(subnettr("1", "0", 2)),10)
	if bcast != "255" {
	  t.Errorf("Invalid broadcast return value for subnettr. Got %s, expected %s.", bcast, "255")
	}
	if subnet != "192" {
	  t.Errorf("Invalid host return value for subnettr. Got %s, expected %s.", subnet, "192")
	}

}

func TestCidr_To_Mask(t *testing.T) {
	nMask := cidr_to_mask("23")
	if nMask != "255.255.254.0" {
	  t.Errorf("Invalid netmask return value for subnettr. Got %s, expected %s.", nMask, "255.255.254.0")
	}

}

func TestSubnettrCore(t *testing.T) {
	network := subnettrCore("192.168.1.5", "27")
	if network[0] != "192.168.1.0" {
	  t.Errorf("Invalid value for network address. Got %s, expected %s.", network[0], "192.168.1.0")
	}
	if network[1] != "192.168.1.30" {
	  t.Errorf("Invalid address for last host on network. Got %s, expected %s.", network[1], "192.168.1.30")
	}
	if network[2] != "192.168.1.1" {
	  t.Errorf("Invalid address for first host on network. Got %s, expected %s.", network[2], "192.168.1.1")
	}
	if network[3] != "192.168.1.31" {
	  t.Errorf("Invalid broadcast address. Got %s, expected %s.", network[3], "192.168.1.31")
	}
	if network[4] != "255.255.255.224" {
	  t.Errorf("Invalid netmask address. Got %s, expected %s.", network[4], "255.255.255.224")
	}
}
