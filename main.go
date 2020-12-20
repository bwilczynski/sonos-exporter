package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var (
	zpIPAddress = flag.String("connect", "", "`zone_player` to connect to")
)

type (
	NetworkInfo struct {
		Items []SupportInfo `xml:"ZPSupportInfo"`
	}

	SupportInfo struct {
		Info     ZonePlayerInfo `xml:"ZPInfo"`
		Commands []Command      `xml:"Command"`
		Files    []File         `xml:"File"`
	}

	ZonePlayerInfo struct {
		Name            string `xml:"ZoneName"`
		SerialNumber    string
		SoftwareVersion string
		HardwareVersion string
		IPAddress       string
		MACAddress      string
	}

	File struct {
		Name string `xml:"name,attr"`
		Text string `xml:",chardata"`
	}

	Command struct {
		Commandline string `xml:"cmdline,attr"`
		Text        string `xml:",chardata"`
	}
)

func (i *SupportInfo) FindCommand(commandLine string) *Command {
	for _, c := range i.Commands {
		if c.Commandline == commandLine {
			return &c
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *zpIPAddress == "" {
		flag.Usage()
		os.Exit(1)
	}

	info, err := fetchNetworkInfo(*zpIPAddress)
	if err != nil {
		panic(err)
	}

	for _, item := range info.Items {
		fmt.Println(item.Info.Name)
		c := item.FindCommand("/sbin/ifconfig")
		fmt.Println(c.Text)
	}
}

func fetchNetworkInfo(zpIPAddress string) (info NetworkInfo, err error) {
	url := fmt.Sprintf("http://%s:1400/support/review", zpIPAddress)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	d := xml.NewDecoder(resp.Body)
	d.Decode(&info)

	return
}
