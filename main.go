package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func checkDomain(domain string) {
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	maxRecords, err := net.LookupMX(domain)

	if err != nil {
		panic(err)
	}

	if len(maxRecords) > 0 {
		hasMx = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		panic(err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		panic(err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%s, %t, %t, %s, %t, %s\n", domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
