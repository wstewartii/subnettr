package main

import "fmt"
import "strconv"
import "os"
import "net"
import "net/http"
import "log"
import "encoding/json"
import "flag"
import "strings"
import "math"
import "math/bits"
import "errors"

type NetworkObject struct {
	NetworkID           net.IP
	SubnetMask          net.IP
	BroadcastAddress    net.IP
	UsableHostAddresses float64
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

func getNetworkObject(addr string, sbnet string) (NetworkObject, error) {

	var netAddress net.IP
	var broadcastAddress net.IP
	var subnet string

	// Check to see if addr can be parsed as an ip address
	if net.ParseIP(addr) == nil {
		return NetworkObject{}, errors.New("Invalid ip address format\n")
	}

	// Check to see if sbnet is a valid subnet mask
	if net.ParseIP(sbnet) == nil {
		// If sbnet is not a subnet mask, check to see if it is a number
		if _, err := strconv.ParseInt(sbnet, 10, 64); err == nil {
			subnet = cidrToMask(sbnet)
		} else {
			return NetworkObject{}, errors.New("Invalid subnet mask/cidr format\n")
		}
	} else {
		subnet = sbnet
	}

	// get network address
	ipAddr := net.ParseIP(addr)
	netMask := net.ParseIP(subnet)
	hostBits := 0
	for i, v := range ipAddr {
		netAddress = append(netAddress, v&netMask[i])
		// invert the last 4 bits in the 16 element array to calculate the broadcast address
		if i > 11 {
			netMaskInverse := ^netMask[i]
			broadcastAddress = append(broadcastAddress, netMaskInverse|v)

			//get number of host bits in the netmask
			hostBits += bits.OnesCount(uint(netMask[i]))
		}
	}

	// Get the number of host bits or 0s in the subnet mask
	numberOfZeros := 32 - hostBits
	numberOfHosts := math.Pow(2, float64(numberOfZeros)) - 2

	netObj := NetworkObject{netAddress, netMask, broadcastAddress, numberOfHosts}

	return netObj, nil
}

func apiUsage(w http.ResponseWriter, r *http.Request) {
	msg := "How to use.\n\n/192.168.1.10/255.255.255.0\n\nor\n\n/172.16.32.22/23\n\n"
	fmt.Fprintf(w, msg)
}

func handleSubnetting(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.URL.Path, "/")) < 3 {
		apiUsage(w, r)
		return
	}

	addr := strings.Split(r.URL.Path, "/")[1]
	sbnet := strings.Split(r.URL.Path, "/")[2]

	resp, err := getNetworkObject(addr, sbnet)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	netObj, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "%s\n", string(netObj))

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
		resp, err := getNetworkObject(addr, sbnet)
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
		http.HandleFunc("/", handleSubnetting)
		fmt.Printf("starting web server on port %s\n", *webPort)
		log.Fatal(http.ListenAndServe(":"+*webPort, nil))
	}
}
