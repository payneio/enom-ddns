package main

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	enom  = `https://dynamic.name-services.com/interface.asp`
	ipURI = `http://checkip.amazonaws.com`
)

var (
	domain, username, password string
)

func die(err error) {
	fmt.Println(err)
	os.Exit(-1)
}

func main() {

	domain = os.Getenv(`DDNS_DOMAIN`)
	username = os.Getenv(`ENOM_UN`)
	password = os.Getenv(`ENOM_PW`)

	if domain == "" || username == "" || password == "" {
		die(errors.New(`set DDNS_DOMAIN, ENOM_PW env vars`))
	}

	d := strings.Split(domain, ".")
	if len(d) != 3 {
		die(errors.New(`only <record>.<sld>.<tld> domains are supported`))
	}
	name, sld, tld := d[0], d[1], d[2]
	zone := fmt.Sprintf(`%s.%s`, sld, tld)

	// Get my IP
	ip, err := GetIP()
	if err != nil {
		die(err)
	}

	// Send Dynamic DNS update
	err = EnomDDNSUpdate(name, zone, ip, username, password)
	if err != nil {
		die(err)
	}

	fmt.Printf("Dynamic DNS updated. %s.%s = %s\n", name, zone, ip)
}

// GetIP uses an IP service at AWS that returns the IP address specified
// in HTTP client's request. In a NAT environment, this will generally be
// the public IP of the WAN router.
func GetIP() (string, error) {
	resp, err := http.Get(ipURI)
	if err != nil {
		return "", err
	}
	ipb, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	ip := strings.TrimSpace(string(ipb))
	return string(ip), err
}

// EnomDDNSUpdate sends an update request to Enom's Dynamic DNS service
func EnomDDNSUpdate(hostname, zone, ipAddress, username, domainPassword string) error {

	// Enom's certificate is invalid, so let's turn of the verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	resp, err := http.Get(
		fmt.Sprintf(enom+
			`?ResponseType=xml`+
			`&Command=SetDNSHost`+
			`&HostName=%s`+
			`&Zone=%s`+
			`&Address=%s`+
			`&UID=%s`+
			`&DomainPassword=%s`,
			hostname,
			zone,
			ipAddress,
			username,
			domainPassword,
		),
	)

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// Parse response
	hr := &CommandResult{}
	err = xml.Unmarshal(body, hr)
	if err != nil {
		return err
	}

	if hr.ErrCount > 0 {
		return errors.New(hr.Errors.Err1)
	}

	return err
}

type Error struct {
	Err1 string `xml:"Err1"`
}

type CommandResult struct {
	Command  string `xml:"Command"`
	Language string `xml:"Language"`
	ErrCount int    `xml:"ErrCount"`
	Errors   Error  `xml:"errors"`
}
