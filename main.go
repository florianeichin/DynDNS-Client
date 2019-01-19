package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	ddns "github.com/jayschwa/go-dyndns"
)

type Configuration struct {
	Url      string
	Username string
	Password string
	Hostname string
}

func main() {
	configuration := Configuration{}
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	s := ddns.Service{
		configuration.Url,
		configuration.Username,
		configuration.Password}
	ip, err := getIP()
	if err != nil {
		panic(err)
	}
	_, err = s.Update(configuration.Hostname, ip)
	if err != nil && strings.Contains(err.Error(), "nochg") {
		panic(err)
	}
}

func getIP() (net.IP, error) {
	response, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	parts := strings.Split(string(content), ".")
	var ips [4]byte
	for i, a := range parts {
		tmp, _ := strconv.Atoi(a)
		ips[i] = byte(tmp)
	}
	return net.IPv4(ips[0], ips[1], ips[2], ips[3]), nil
}
