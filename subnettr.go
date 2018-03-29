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
import "regexp"
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
	hostBits := 0

	maskFormat, err := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", sbnet)
	if err != nil {
		return NetworkObject{}, err
	}
	cidrFormat, err := regexp.MatchString("^[0-9]{1,2}$", sbnet)
	if err != nil {
		return NetworkObject{}, err
	}

	if maskFormat == false {
		if cidrFormat == true {
			subnet = cidrToMask(sbnet)
		} else {
			return NetworkObject{}, errors.New("Error: invalid netmask format!\n")
		}
	} else {
		subnet = sbnet
	}

	// get network address
	ipAddr := net.ParseIP(addr)
	netMask := net.ParseIP(subnet)
	for i, v := range ipAddr {
		netAddress = append(netAddress, v&netMask[i])
		// invert the last 4 bytes in the array to calculate the broadcast address
		if i > 11 {
			netMaskInverse := ^netMask[i]
			broadcastAddress = append(broadcastAddress, netMaskInverse|v)

			//get number of hosts
			hostBits += bits.OnesCount(uint(netMask[i]))
		}
	}

	numberOfZeros := 32 - hostBits
	numberOfHosts := math.Pow(2, float64(numberOfZeros)) - 2
	netObj := NetworkObject{netAddress, netMask, broadcastAddress, numberOfHosts}

	return netObj, nil
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
		http.HandleFunc("/subnet/", handleSubnetting)
		fmt.Printf("starting web server on port %s\n", *webPort)
		log.Fatal(http.ListenAndServe(":"+*webPort, nil))
	}
}
