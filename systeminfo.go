package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
	"time"

	"github.com/mackerelio/go-osstat/uptime"
)

func getSystemUptime() (time.Duration) {
	uptime, err := uptime.Get()
	if err != nil {
		fmt.Printf("%s\n", err)
		return time.Duration(0)
	}
	return time.Duration(uptime)
}

func getUname() (domainname string, machine string, nodename string, release string, sysname string, version string) {
	var uts syscall.Utsname
	syscall.Uname(&uts)
	
	tmp := ""
	for _, r := range uts.Nodename {
		if r == 0 {
			break
		}
		tmp += string(r)
	}
	nodename = tmp

	tmp = ""
	for _, r := range uts.Domainname {
		if r == 0 {
			break
		}
		tmp += string(r)
	}
	domainname = tmp

	tmp = ""
	for _, r := range uts.Machine {
		if r == 0 {
			break
		}
		tmp += string(r)
	}
	machine = tmp

	tmp = ""
	for _, r := range uts.Release {
		if r == 0 {
			break
		}
		tmp += string(r)
	}
	release = tmp

	tmp = ""
	for _, r := range uts.Sysname {
		if r == 0 {
			break
		}
		tmp += string(r)
	}
	sysname = tmp

	tmp = ""
	for _, r := range uts.Version {
		if r == 0 {
			break
		}
		tmp += string(r)
	}
	version = tmp

	return
}

func getDistroName() (name string) {
	val, err := ioutil.ReadFile("/etc/issue")
	if (err != nil) {
		val = []byte("Unknown")
	}
	name = string(val)
	name = strings.Replace(name, "\\n", "", -1)
	name = strings.Replace(name, "\\l", "", -1)
	return
}