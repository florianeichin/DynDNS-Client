package main

import (
	"encoding/json"
	"fmt"
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
	configuration := []Configuration{}
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	ip, err := getGlobalIP()
	if err != nil {
		panic(err)
	}

	for _, config := range configuration {
		s := ddns.Service{
			config.Url,
			config.Username,
			config.Password}
		currentIP, err := net.LookupIP(config.Hostname)
		if err != nil {
			fmt.Printf("[ERROR] Lookup failed: %s \n", config.Hostname)
			continue
		}
		if contains(currentIP, ip) {
			fmt.Printf("[INFO] nothing changed: %s \n", config.Hostname)
			continue
		}
		_, err = s.Update(config.Hostname, ip)
		if err == nil {
			fmt.Printf("[INFO] updated: %s \n", config.Hostname)
		} else {
			panic(err)
		}
	}
}

func getGlobalIP() (net.IP, error) {
	response, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(content), ".")
	var ips [4]byte
	for i, a := range parts {
		tmp, _ := strconv.Atoi(a)
		ips[i] = byte(tmp)
	}
	return net.IPv4(ips[0], ips[1], ips[2], ips[3]), nil
}

func contains(currentIPS []net.IP, ip net.IP) bool {
	for _, i := range currentIPS {
		if i.Equal(ip) {
			return true
		}
	}
	return false
}
