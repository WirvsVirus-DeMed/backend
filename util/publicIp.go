package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func GetPublicIp() net.IP {
	url := "https://ident.me/"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[*] External Ip:%s\n", ip)
	return net.ParseIP(string(ip))
}
