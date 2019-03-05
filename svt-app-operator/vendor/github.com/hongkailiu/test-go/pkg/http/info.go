package http

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// VERSION of the http cmd
	VERSION = "2.0.0"
)

//Info represents information
type Info struct {
	Version string    `json:"version"`
	Ips     []string  `json:"ips"`
	Now     time.Time `json:"now"`
}

// getInfo returns the required information
func getInfo() *Info {

	i := Info{}
	i.Version = VERSION
	i.Ips = getIps()
	i.Now = time.Now()
	return &i
}

func getIps() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error(err.Error())
		return []string{err.Error()}
	}
	// handle err
	var result []string
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			result = append(result, err.Error())
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			result = append(result, ip.String())
		}
	}
	return result
}
