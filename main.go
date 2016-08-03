package main

import (
	"flag"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

var (
	dnsCache = make(map[string]string)
	emails   = flag.String("addrs", "", "List of email addresses separeted by commas.")
)

func main() {
	flag.Parse()

	addrs := []string{}
	for _, e := range strings.Split(*emails, ",") {
		trimmed := strings.TrimSpace(e)
		if trimmed != "" {
			addrs = append(addrs, trimmed)
		}
	}
	for _, addr := range addrs {
		ok, err := accountExists(addr)
		if err != nil {
			log("[error] " + addr + ": " + err.Error())
		}
		logData(fmt.Sprintf("%v, %v", addr, ok))
	}
}

func accountExists(email string) (bool, error) {
	server, err := findMailserver(email)
	if err != nil {
		return false, err
	}
	log("smtp: connecting", server)
	conn, err := smtp.Dial(server + ":25")
	if err != nil {
		return false, err
	}
	log("smtp: connect")
	defer conn.Quit()

	if err := conn.Mail("sender@example.com"); err != nil {
		return false, err
	}
	log("smtp: from")

	if err := conn.Rcpt(email); err != nil {
		return false, nil
	}
	log("smtp: to")
	return true, nil
}

func findMailserver(rawEmail string) (string, error) {
	email, err := mail.ParseAddress(rawEmail)
	if err != nil {
		return "", err
	}

	// Get hostname from email address.
	idx := strings.Index(email.Address, "@")
	if idx < 0 {
		return "", fmt.Errorf("cant parse email: %v", email.Address)
	}
	host := email.Address[idx+1 : len(email.Address)]

	if dnsCache[host] != "" {
		return dnsCache[host], nil
	}

	// Lookup DNS MX record
	mxs, err := net.LookupMX(host)
	if err != nil {
		return "", err
	}

	if len(mxs) == 0 {
		return "", fmt.Errorf("%v don't have MX server", host)
	}
	dnsCache[host] = mxs[0].Host
	return mxs[0].Host, nil
}

func logData(s ...interface{}) {
	fmt.Fprintln(os.Stdout, s...)
}

func log(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
}
