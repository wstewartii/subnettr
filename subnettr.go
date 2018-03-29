package main

import "fmt"
import "strconv"
import "errors"
import "os"
import "net/http"
import "log"
import "encoding/json"
import "flag"
import "strings"
import "regexp"

type NetworkInfo struct {
	FirstHostAddress string
	LastHostAddress  string
	Subnet           string
	BroadcastAddress string
	SubnetMask       string
}

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

func cidrToMask(cidr string) string {

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

func subnettrCore(addr string, sbnet string) (NetworkInfo, error) {

	var nmask string
	var sbnetList []string
	var bcastList []string
	var lhostList []string
	var fhostList []string

	addrFormat, aerr := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", addr)
	if aerr != nil {
		return NetworkInfo{}, aerr
	}
	maskFormat, merr := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", sbnet)
	if merr != nil {
		return NetworkInfo{}, merr
	}
	cidrFormat, cerr := regexp.MatchString("^[0-9]{1,2}$", sbnet)
	if cerr != nil {
		return NetworkInfo{}, cerr
	}

	if addrFormat == false {
		return NetworkInfo{}, errors.New("Error: invalid address format!\n")
	}
	if maskFormat == true {
		nmask = sbnet
	} else if cidrFormat == true {
		nmask = cidrToMask(sbnet)
	} else {
		return NetworkInfo{}, errors.New("Error: invalid netmask format!\n")
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

	netInfo := NetworkInfo{firstHost, lastHost, subnet, broadcast, netmask}

	return netInfo, nil
}

func apiUsage(w http.ResponseWriter, r *http.Request) {
	msg := "How to use.\n\n/subnet/192.168.1.10/255.255.255.0\n\nor\n\n/subnet/172.16.32.22/23\n\n"
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

	nInfo, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "%s\n", string(nInfo))

}

func main() {

	webPort := flag.String("port", "8080", "web server port")

	webServer := flag.Bool("server", false, "start a web server")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: subnettr <ip address> <subnet mask or cidr>\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	_, err := strconv.ParseFloat(*webPort, 64)

	if err != nil {
		fmt.Printf("%s is not a valid port number\n", *webPort)
		os.Exit(0)
	}

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

		nInfo, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(nInfo))

	} else {
		http.HandleFunc("/subnet/", handleSubnetting)
		fmt.Printf("starting web server on port %s\n", *webPort)
		log.Fatal(http.ListenAndServe(":"+*webPort, nil))
	}
}
