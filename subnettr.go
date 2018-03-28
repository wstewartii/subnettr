package main

import "fmt"
import "strconv"
import "errors"
import "os"
import "net/http"
import "log"

import "flag"
import "strings"
import "regexp"

func subnettr(addr string, sbnet string, query int) string {

	var resp uint8
	var conv_resp string

	address, aErr := strconv.ParseUint(addr, 10, 8)
	if aErr != nil {
		fmt.Println(aErr)
	}

	subnet, sErr := strconv.ParseUint(sbnet, 10, 64)
	if sErr != nil {
		fmt.Println(sErr)
	}

	netMask := uint8(subnet)
	ipAddr := uint8(address)
	netAddress := netMask & ipAddr
	inverse := ^netMask
	broadcastAddress := netAddress | inverse

	switch query {
	case 1:
		resp = netAddress
		_ = resp
	case 2:
		resp = broadcastAddress
		_ = resp
	case 3:
		resp = netAddress + 1
		_ = resp
	case 4:
		resp = broadcastAddress - 1
		_ = resp
	default:
		resp = netMask
		_ = resp
	}

	conv_resp = strconv.FormatUint(uint64(resp), 10)
	return conv_resp
}

func cidr_to_mask(cidr string) string {

	var maskList []string
	var netMask string

	cidrInt, err := strconv.ParseUint(cidr, 10, 8)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 4; i++ {
		tmpstring := ""
		for ii := 0; ii < 8; ii++ {
			if cidrInt > 0 {
				tmpstring = tmpstring + "1"
				cidrInt--
			} else {
				tmpstring = tmpstring + "0"
			}
		}
		tmpint, ierr := strconv.ParseUint(tmpstring, 2, 64)
		if ierr != nil {
			fmt.Println(ierr)
		}
		maskList = append(maskList, strconv.FormatUint(uint64(tmpint), 10))
	}
	netMask = strings.Join(maskList, ".")

	return netMask
}

func subnettrCore(addr string, sbnet string) ([]string, error) {

	var respList []string
	var nmask string
	var sbnetList []string
	var bcastList []string
	var lhostList []string
	var fhostList []string

	addrFormat, aerr := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", addr)
	if aerr != nil {
		return nil, aerr
	}
	maskFormat, merr := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", sbnet)
	if merr != nil {
		return nil, merr
	}
	cidrFormat, cerr := regexp.MatchString("^[0-9]{1,2}$", sbnet)
	if cerr != nil {
		return nil, cerr
	}

	if addrFormat == false {
		return nil, errors.New("Error: invalid address format!\n")
	}
	if maskFormat == true {
		nmask = sbnet
	} else if cidrFormat == true {
		nmask = cidr_to_mask(sbnet)
	} else {
		return nil, errors.New("Error: invalid netmask format!\n")
	}

	addrList := strings.Split(addr, ".")
	nmaskList := strings.Split(nmask, ".")

	for i, v := range addrList {
		sbnetList = append(sbnetList, subnettr(v, nmaskList[i], 1))
	}
	for i, v := range addrList {
		bcastList = append(bcastList, subnettr(v, nmaskList[i], 2))
	}
	for i, v := range addrList {
		if i == 3 {
			fhostList = append(fhostList, subnettr(v, nmaskList[i], 3))
		} else {
			fhostList = append(fhostList, subnettr(v, nmaskList[i], 1))

		}
	}
	for i, v := range addrList {
		if i == 3 {
			lhostList = append(lhostList, subnettr(v, nmaskList[i], 4))
		} else {
			lhostList = append(lhostList, subnettr(v, nmaskList[i], 2))

		}
	}

	subnet := strings.Join(sbnetList, ".")
	lastHost := strings.Join(lhostList, ".")
	firstHost := strings.Join(fhostList, ".")
	broadcast := strings.Join(bcastList, ".")
	netmask := strings.Join(nmaskList, ".")

	respList = append(respList, subnet)
	respList = append(respList, lastHost)
	respList = append(respList, firstHost)
	respList = append(respList, broadcast)
	respList = append(respList, netmask)

	return respList, nil
}

func apiUsage(w http.ResponseWriter, r *http.Request) {
	msg := "How to use.\n\n/subnet/192.168.1.10/255.255.255.0\n\n"
	fmt.Fprintf(w, msg)
}

func handleSubnetting(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.URL.Path, "/")) < 4 {
		apiUsage(w, r)
		return
	}

	addr := strings.Split(r.URL.Path, "/")[2]
	sbnet := strings.Split(r.URL.Path, "/")[3]

	resp, err := subnettrCore(addr, sbnet)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	msg := "Address Range: " + resp[2] + "-" + resp[1] + "\n"
	msg += "Net Address: " + resp[0] + "\n"
	msg += "Broadcast Address: " + resp[3] + "\n"
	msg += "Subnet Mask: " + resp[4] + "\n"
	fmt.Fprintf(w, msg)

}

func main() {

	webPort := "8080"

	webServer := flag.Bool("server", false, "start a web server")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: subnettr <ip address> <subnet mask>\n\n<ip address> The ip address e.g., 192.168.0.1\n\n<subnet mask> The subnet mask or cidr block e.g., 255.255.255.0 or 24\n\noptions:\n\n\t-server	start a web server\n\n")
	}
	flag.Parse()
	if *webServer == false {
		if len(flag.Args()) < 2 && *webServer == false {
			flag.Usage()
			os.Exit(0)
		}
		addr := flag.Args()[0]
		sbnet := flag.Args()[1]
		resp, err := subnettrCore(addr, sbnet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		subnet := resp[0]
		lastHost := resp[1]
		firstHost := resp[2]
		broadcast := resp[3]
		netmask := resp[4]

		fmt.Println("Address Range: " + firstHost + "-" + lastHost)
		fmt.Println("Net Address: " + subnet)
		fmt.Println("Broadcast Address: " + broadcast)
		fmt.Println("Subnet Mask: " + netmask)
	} else {
		http.HandleFunc("/subnet/", handleSubnetting)
		fmt.Println("starting web server on port", webPort)
		log.Fatal(http.ListenAndServe(":"+webPort, nil))
	}
}
