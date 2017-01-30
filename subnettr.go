package main

import "fmt"
import "strconv"
import "os"
import "strings"
import "regexp"

func subnettr(addr string, sbnet string, query int) uint8 {

var resp uint8

address, aErr := strconv.ParseUint(addr, 10, 8);
if aErr != nil {
        fmt.Println(aErr)
  }

subnet, sErr := strconv.ParseUint(sbnet, 10, 64);
if sErr != nil {
        fmt.Println(sErr)
    }

netMask := uint8(subnet)
ipAddr := uint8(address)
netAddress := netMask&ipAddr
inverse := ^netMask
broadcastAddress := netAddress|inverse

switch query {
case 1:
  resp = netAddress
  _ = resp
case 2:
  resp = broadcastAddress
  _ = resp
default:
  resp = netMask
  _ = resp
}
return resp
}

func cidr_to_mask(cidr string) string {

var maskList []string
var netMask string

cidrInt, err := strconv.ParseUint(cidr, 10, 8);
if err != nil {
fmt.Println(err)
}

for i :=0;i<4;i++ {
  tmpstring := ""
  for ii :=0;ii<8;ii++ {
  if cidrInt > 0 {
  tmpstring = tmpstring + "1"
  cidrInt --
  } else {
  tmpstring = tmpstring + "0"
  }
  }
  tmpint, ierr := strconv.ParseUint(tmpstring, 2, 64)
  if ierr != nil {
    fmt.Println(ierr)
  }
  maskList = append(maskList, strconv.FormatUint(uint64(tmpint),10))
}
netMask = strings.Join(maskList, ".")

return netMask
}

func main() {

var nmask string
var sbnetList []string
var bcastList []string
var lhostList []string

if len(os.Args) < 2 {
  fmt.Println("Error: address and subnet is required")
  os.Exit(1)
}

maskFormat, merr := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$",os.Args[2])
if merr != nil {
  fmt.Println(merr)
}
cidrFormat, cerr := regexp.MatchString("^[0-9]{2,2}$",os.Args[2])
if cerr != nil {
  fmt.Println(cerr)
}

if maskFormat == true {
nmask = os.Args[2];
} else if cidrFormat == true {
nmask = cidr_to_mask(os.Args[2]);
} else {
fmt.Println("Error: invalid netmask format!")
os.Exit(1)
}

addr := os.Args[1];

addrList := strings.Split(addr, ".")
nmaskList := strings.Split(nmask, ".")

for i,v := range addrList {
  sbnetList = append(sbnetList,strconv.FormatUint(uint64(subnettr(v, nmaskList[i], 1)),10))
}
for i,v := range addrList {
  bcastList = append(bcastList,strconv.FormatUint(uint64(subnettr(v, nmaskList[i], 2)),10))
}
for i,v := range addrList {
  if i == 3 {
  lhostList = append(lhostList,strconv.FormatUint(uint64(subnettr(v, nmaskList[i], 2)-1),10))
  } else {
  lhostList = append(lhostList,strconv.FormatUint(uint64(subnettr(v, nmaskList[i], 2)),10))

  }
}

subnet := strings.Join(sbnetList, ".")
lastHost := strings.Join(lhostList, ".")
broadcast := strings.Join(bcastList, ".")

fmt.Println("Net Address: " + subnet)
fmt.Println("Broadcast Address: " + broadcast)
fmt.Println("Subnet Range: " + subnet + "-" + lastHost)


}
