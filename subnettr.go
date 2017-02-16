package main

import "fmt"
import "strconv"
import "os"
import "flag"
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
var fhostList []string


flag.Usage = func() {
    fmt.Fprintf(os.Stdout, "Usage: subnettr <ip address> <subnet mask>\n\n<ip address> The ip address e.g., 192.168.0.1\n\n<subnet mask> The subnet mask or cidr block e.g., 255.255.255.0 or 24\n\n")
}
flag.Parse()

if len(flag.Args()) < 2 {
  flag.Usage()
  os.Exit(1)
}

addrFormat, aerr := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$",flag.Args()[0])
if aerr != nil {
  fmt.Println(aerr)
}
maskFormat, merr := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$",flag.Args()[1])
if merr != nil {
  fmt.Println(merr)
}
cidrFormat, cerr := regexp.MatchString("^[0-9]{1,2}$",flag.Args()[1])
if cerr != nil {
  fmt.Println(cerr)
}

if addrFormat == false {
fmt.Println("Error: invalid address format!")
os.Exit(1)
}
if maskFormat == true {
nmask = flag.Args()[1];
} else if cidrFormat == true {
nmask = cidr_to_mask(flag.Args()[1]);
} else {
fmt.Println("Error: invalid netmask format!")
os.Exit(1)
}

addr := flag.Args()[0];

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
  fhostList = append(fhostList,strconv.FormatUint(uint64(subnettr(v, nmaskList[i], 1)+1),10))
  } else {
  fhostList = append(fhostList,strconv.FormatUint(uint64(subnettr(v, nmaskList[i], 1)),10))

  }
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
firstHost := strings.Join(fhostList, ".")
broadcast := strings.Join(bcastList, ".")
netmask := strings.Join(nmaskList, ".")

fmt.Println("Address Range: " + firstHost + "-" + lastHost)
fmt.Println("Net Address: " + subnet)
fmt.Println("Broadcast Address: " + broadcast)
fmt.Println("Subnet Mask: " + netmask)


}
