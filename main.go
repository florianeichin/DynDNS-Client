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
	log "github.com/sirupsen/logrus"
)

// Configuration holds the complete JSON configuration data
type Configuration struct {
	Logfile   string
	DNSConfig []DNSConfig
}

// DNSConfig holds all data, that is used for connecting to the DNS Service
type DNSConfig struct {
	URL      string
	Username string
	Password string
	Hostname string
}

func main() {
	configuration := Configuration{}
	file, err := os.Open("./config.json")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}

	logfile, err := os.OpenFile(configuration.Logfile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	log.SetOutput(logfile)

	ip, err := getGlobalIP()
	if err != nil {
		log.Panic(err)
	}

	for _, config := range configuration.DNSConfig {
		s := ddns.Service{URL: config.URL, Username: config.Username, Password: config.Password}
		currentIP, err := net.LookupIP(config.Hostname)
		if err != nil {
			log.Errorf("Lookup failed: %s", config.Hostname)
			continue
		}
		if contains(currentIP, ip) {
			log.Infof("nothing changed: %s", config.Hostname)
			continue
		}
		_, err = s.Update(config.Hostname, ip)
		if err == nil {
			log.Infof("updated: %s", config.Hostname)
		} else {
			log.Panic(err)
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
