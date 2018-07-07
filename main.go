package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	enom  = `https://reseller.enom.com/interface.asp`
	ipURI = `http://checkip.amazonaws.com`
)

var domain, username, password string

func die(err error) {
	fmt.Println(err)
	os.Exit(-1)
}

func main() {

	domain = os.Getenv(`DDNS_DOMAIN`)
	username = os.Getenv(`DDNS_UN`)
	password = os.Getenv(`DDNS_PW`)

	if domain == "" || username == "" || password == "" {
		die(errors.New(`set DDNS_DOMAIN, DDNS_UN DDNS_PW env vars`))
	}

	d := strings.Split(domain, ".")
	if len(d) != 3 {
		die(errors.New(`only <record>.<sld>.<tld> domains are supported`))
	}
	name, sld, tld := d[0], d[1], d[2]

	// Get my IP
	ip, err := GetIP()
	if err != nil {
		die(err)
	}

	// Get records for the domain
	records, err := GetRecords(tld, sld)
	if err != nil {
		die(err)
	}

	// Update the records if necessary
	var (
		dirty  bool
		found  bool
		update []Record
	)
	for _, record := range records {
		h := Record{
			Name:    record.Name,
			Type:    record.Type,
			Address: record.Address,
			MXPref:  record.MXPref,
		}
		if record.Name == name {
			found = true
			if record.Address != ip {
				dirty = true
				h.Address = ip
				h.Type = `A`
			}
			if record.Name == "" {
				record.Name = `@`
			}
		}
		update = append(update, h)
	}
	if !found {
		update = append(update, Record{
			Name:    name,
			Type:    `A`,
			Address: ip,
		})
	}

	// Check and/or modify a record here
	if dirty || !found {
		err := SetRecords(domain, update)
		if err != nil {
			die(err)
		}
		fmt.Printf("%s updated to %s\n", domain, ip)
	} else {
		fmt.Printf("%s already set to %s. Skipping update.\n", domain, ip)
	}
}

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

func GetRecords(tld, sld string) ([]Record, error) {

	resp, err := http.PostForm(
		enom,
		url.Values{
			"UID":          []string{username},
			"PW":           []string{password},
			"Command":      []string{`GetHosts`},
			"TLD":          []string{tld},
			"SLD":          []string{sld},
			"ResponseType": []string{`XML`},
		},
	)
	if err != nil {
		return nil, err
	}

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Parse response
	hr := &GetRecordsResponse{}
	err = xml.Unmarshal(body, hr)
	if err != nil {
		return nil, err
	}

	if hr.ErrCount > 0 {
		return nil, errors.New(hr.Errors.Err1)
	}

	return hr.Records, nil
}

func SetRecords(domain string, records []Record) error {

	fragments := strings.Split(domain, ".")
	params := url.Values{
		"UID":          []string{username},
		"PW":           []string{password},
		"Command":      []string{`SetHosts`},
		"TLD":          []string{fragments[1]},
		"SLD":          []string{fragments[0]},
		"ResponseType": []string{`XML`},
	}

	for i, record := range records {
		si := strconv.Itoa(i + 1)
		params.Set(`RecordName`+si, record.Name)
		params.Set(`RecordType`+si, record.Type)
		params.Set(`Address`+si, record.Address)
		if record.MXPref != "" {
			params.Set(`MXPref`+si, record.MXPref)
		}
	}

	resp, err := http.PostForm(enom, params)

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

type Record struct {
	Name     string `xml:"name"`
	Type     string `xml:"type"`
	Address  string `xml:"address"`
	RecordID string `xml:"hostid"`
	MXPref   string `xml:"mxpref"`
}

type GetRecordsResponse struct {
	CommandResult
	Records []Record `xml:"host"`
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
